package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/dto"
	"sca-integrator/app/helper"
	"strings"
)

type projectPathInfo struct {
	scannedProjectFilePath string
	projectPath            string
	typeVersion            string
}

type projectOptionInfo struct {
	unfixedOptions  []entity.ProjectFilterOption
	severityOptions []entity.ProjectFilterOption
}

type argumentInfo struct {
	isUnfixedOption  bool
	severityArgument string
	skipArgument     string
}

func (p ProjectServiceImpl) scanningRepository(ctx *gin.Context, project entity.Project) {
	p.stdLog.InfoFunction(fmt.Sprintf("START SCANNING REPOSITORY %s", project.Key))

	pathData := p.preparingProjectPath(ctx, project)
	p.prepareRunScanAndAction(ctx, project, pathData)
	p.saveResultToDb(ctx, project, pathData)

	p.stdLog.InfoFunction(fmt.Sprintf("END SCANNING REPOSITORY %s", project.Key))
}

// PREPARING PROJECT PATH PHASE
func (p ProjectServiceImpl) preparingProjectPath(ctx *gin.Context, project entity.Project) projectPathInfo {
	var curDir, _ = os.Getwd()
	typeVersion := "repository"

	// Example Result of Project Path => /_project-repository/workspace/test_project-sca
	RepositoryPath := fmt.Sprintf("%s%s", curDir, p.prefixProjectFolder)
	projectPath := fmt.Sprintf("%s%s-sca", RepositoryPath, project.Key)

	// Example Result of Scanned Project File Path => /_scanned-project-files/scan_repo_file_test_project_repository-.json
	ProjectFilePath := fmt.Sprintf("%s%s", curDir, p.prefixResultFolder)
	projectFileName := "scan_repo_file_" + project.Key + "_" + typeVersion + "-sca" + ".json"
	scannedProjectFilePath := ProjectFilePath + projectFileName

	data := projectPathInfo{
		scannedProjectFilePath: scannedProjectFilePath,
		projectPath:            projectPath,
		typeVersion:            typeVersion,
	}

	return data
}

// PREPARING RUN SCAN AND ACTION PHASE
func (p ProjectServiceImpl) prepareRunScanAndAction(ctx *gin.Context, project entity.Project, pathData projectPathInfo) projectOptionInfo {
	var argumentData argumentInfo
	var projectOptionData projectOptionInfo

	argumentData.isUnfixedOption, projectOptionData.unfixedOptions = p.addArgumentUnfixedOption(ctx, project.Id)

	argumentData.severityArgument, projectOptionData.severityOptions = p.addArgumentSeverity(ctx, project.Id)

	argumentData.skipArgument = p.addArgumentSkip(ctx, project.Id, project.Key)

	out, err := p.runScanCommand(argumentData, pathData)

	if err != nil {
		outPutCommand := fmt.Sprintf("%q\n", out)
		p.stdLog.InfoFunction("Error scan command: " + outPutCommand)

		var succesDownloadDB bool
		if strings.Contains(outPutCommand, "--skip-update cannot be specified on the first run") {
			succesDownloadDB = p.trivyCli.DownloadTrivyDb()
		}
		if succesDownloadDB {
			out, err = p.runScanCommand(argumentData, pathData)
		}
		p.checkErrorScan(ctx, project, err)
	}
	p.stdLog.InfoFunction("Success scan command: " + fmt.Sprintf("%q\n", out))

	return projectOptionData
}

func (p ProjectServiceImpl) addArgumentUnfixedOption(ctx *gin.Context, projectId int) (bool, []entity.ProjectFilterOption) {
	result := false
	unfixedOptions := p.repoOption.GetAllByProjectIdAndFilterType(ctx, p.db, projectId, "Hide Unfixed Vulnerabilities")
	if len(unfixedOptions) > 0 {
		result = true
	}
	return result, unfixedOptions
}

func (p ProjectServiceImpl) addArgumentSeverity(ctx *gin.Context, projectId int) (string, []entity.ProjectFilterOption) {
	var result string
	severityOptions := p.repoOption.GetAllByProjectIdAndFilterType(ctx, p.db, projectId, "Severity")
	if len(severityOptions) > 0 {
		var severity string
		length := len(severityOptions)
		for index, severityOption := range severityOptions {
			if length == index+1 {
				severity += fmt.Sprintf("%s%s", severity, severityOption.Value)
			} else {
				severity += fmt.Sprintf("%s%s,", severity, severityOption.Value)
			}
		}
		result += fmt.Sprintf(" --severity %s", severity)
	}
	return result, severityOptions
}

func (p ProjectServiceImpl) addArgumentSkip(ctx *gin.Context, projectId int, projectKey string) string {
	exclusions := p.repoExclusion.GetAllByProjectId(ctx, p.db, projectId)
	var result string
	for _, exclusion := range exclusions {
		exclusion.Path = strings.Replace(exclusion.Path, projectKey+":", "", -1)
		argumentFlag := ""
		if exclusion.Type == "DIR" || exclusion.Type == "TRK" {
			argumentFlag = "--skip-dir"
		} else if exclusion.Type == "FIL" || exclusion.Type == "UTS" {
			argumentFlag = "--skip-files"
		}
		result += fmt.Sprintf(" %s '%s'", argumentFlag, exclusion.Path)
	}
	result += " --skip-dirs 'node_modules/' --skip-dirs 'vendor/' --skip-dirs 'build/' --skip-dirs 'dist/' --skip-dirs 'target/'"

	return result
}

func (p ProjectServiceImpl) runScanCommand(argumentData argumentInfo, pathData projectPathInfo) ([]byte, error) {
	trivyCli := p.trivyCli.Init().AddFileSystemArgument(argumentData.skipArgument)
	if argumentData.isUnfixedOption {
		trivyCli.AddIgnoreUnfixedArgument()
	}
	trivyCli.AddSeverityArgument(argumentData.severityArgument).
		AddSecurityCheckArgument("vuln").
		AddSkipDBUpdateArgument().
		AddFormat("json").
		AddOutput(pathData.scannedProjectFilePath).
		AddProjectPath(pathData.projectPath)
	out, err := trivyCli.Exec()

	return out, err
}

func (p ProjectServiceImpl) checkErrorScan(ctx *gin.Context, project entity.Project, err error) {
	if err != nil {
		p.failedScan(ctx, project)
		helper.ErrorHandler(err)
	}
}

// SAVE RESULT TO DB PHASE
func (p ProjectServiceImpl) saveResultToDb(ctx *gin.Context, project entity.Project, pathData projectPathInfo) {
	scanVersion := p.incrementScanVersion(ctx, project)

	results := p.mappingResultOutput(pathData.scannedProjectFilePath, project.Id, project.Key)

	nothingToChange := p.compareNewResultWithPrevResult(ctx, project.Id, results, scanVersion)

	if nothingToChange {
		project.StatusScan = 3
		//project.Message = "Nothing To Update"
		p.repo.Update(ctx, p.db, project)

		p.stdLog.InfoFunction("Done scanning project repository... (nothing to update)")
		return
	}

	for _, result := range results {
		resultRule := p.repoResult.GetOneResultByProjectIdAndRuleAndTargetFile(ctx, p.db, project.Id, result.Rule, result.TargetFile)
		p.addOrUpdateResult(ctx, result, resultRule, scanVersion, pathData.typeVersion)
	}

	p.deleteAllPrevVersionResult(ctx, project, scanVersion)

	project.StatusScan = 3
	//project.StatusMessage = ""
	project.CurrentScanVersion = scanVersion
	p.repo.Update(ctx, p.db, project)
}

func (p ProjectServiceImpl) incrementScanVersion(ctx *gin.Context, project entity.Project) int {
	lastScanResult := p.repoResult.GetLastByProjectId(ctx, p.db, project.Id)

	scanVersion := 0

	if lastScanResult.Id > 0 {
		if lastScanResult.ScanVersion <= project.CurrentScanVersion {
			scanVersion = project.CurrentScanVersion + 1
		} else {
			scanVersion = int(lastScanResult.ScanVersion) + 1
		}
	} else {
		scanVersion++
	}

	return scanVersion
}

func (p ProjectServiceImpl) mappingResultOutput(filePathResult string, projectId int, projectKey string) []entity.Result {
	var data dto.ResultOutputFile

	sourceFile, _ := os.Open(filePathResult)
	defer sourceFile.Close()

	byteValue1, _ := ioutil.ReadAll(sourceFile)
	_ = json.Unmarshal(byteValue1, &data)

	results := []entity.Result{}

	for _, result := range data.Results {
		for _, vulnerability := range result.Vulnerabilities {
			tmpResult := new(entity.Result)
			tmpResult.TargetFile = result.Target
			tmpResult.PackagesType = result.Type
			tmpResult.ProjectId = projectId
			tmpResult.ProjectKey = projectKey
			tmpResult.PrimaryUrl = vulnerability.PrimaryURL
			tmpResult.Rule = vulnerability.VulnerabilityID
			tmpResult.PackageName = vulnerability.PkgName
			tmpResult.InstalledVersion = vulnerability.InstalledVersion
			tmpResult.FixedVersion = vulnerability.FixedVersion
			tmpResult.References = strings.Join(vulnerability.References, "|")
			tmpResult.Title = vulnerability.Title
			tmpResult.Description = vulnerability.Description
			tmpResult.Severity = vulnerability.Severity
			tmpResult.PublishedDate = vulnerability.PublishedDate
			tmpResult.LastModifiedDate = vulnerability.LastModifiedDate
			tmpResult.CvssSource = "nvd"
			tmpResult.CvssV3 = fmt.Sprintf("%f", vulnerability.Cvss.Nvd.V3Score) + "|" + vulnerability.Cvss.Nvd.V3Vector
			tmpResult.CvssV2 = fmt.Sprintf("%f", vulnerability.Cvss.Nvd.V3Score) + "|" + vulnerability.Cvss.Nvd.V2Vector

			results = append(results, *tmpResult)
		}
	}
	return results
}

func (p ProjectServiceImpl) compareNewResultWithPrevResult(ctx *gin.Context, projectId int, newResults []entity.Result, currentVersion int) bool {
	var fixedVulnerabilityBool = false
	var newVulnerabilityBool = false

	prevResults := p.repoResult.GetAllByProjectIdAndScanVersion(ctx, p.db, projectId, currentVersion-1)
	newVulnerabilityBool = p.checkingNotHaveSameRule(prevResults, newResults)
	fixedVulnerabilityBool = p.checkingNotHaveSameRule(prevResults, prevResults)

	if fixedVulnerabilityBool == false && newVulnerabilityBool == false {
		return true
	}
	return false
}

func (p ProjectServiceImpl) checkingNotHaveSameRule(result1 []entity.Result, result2 []entity.Result) bool {
	if len(result1) == 0 {
		return true
	}

	ruleSet := make(map[string]bool)

	for _, result := range result1 {
		if _, ok := ruleSet[result.Rule]; !ok {
			ruleSet[result.Rule] = true
		}
	}

	for _, result := range result2 {
		if _, ok := ruleSet[result.Rule]; ok {
			continue
		}
		return true
	}
	return false
}

func (p ProjectServiceImpl) addOrUpdateResult(ctx *gin.Context, result entity.Result, oldResult entity.Result, scanVersion int, typeVersion string) {
	if oldResult.Id != 0 {
		tmpScanVersion := result.ScanVersion
		tmpLastUpdate := result.UpdatedAt.Format("2006-01-02 15:04:05")

		result.ScanVersion = scanVersion
		result.StatusResult = 0
		result.LastFoundAt = fmt.Sprintf("Scan no. %d at %s", tmpScanVersion, tmpLastUpdate)
		p.repoResult.Update(ctx, p.db, result)
	} else {
		result.ScanType = typeVersion
		result.ScanVersion = scanVersion
		result.LastFoundAt = ""
		p.repoResult.Create(ctx, p.db, result)
	}
}

func (p ProjectServiceImpl) deleteAllPrevVersionResult(ctx *gin.Context, project entity.Project, scanVersion int) {
	prevResults := p.repoResult.GetAllByProjectIdAndScanVersion(ctx, p.db, project.Id, scanVersion-1)
	for _, v := range prevResults {
		p.repoResult.DeleteOne(ctx, p.db, v)
	}
}

// UTILS SECTION
func (p ProjectServiceImpl) failedScan(ctx *gin.Context, project entity.Project) {
	project.StatusScan = 2
	// projectInformation.StatusMessage = statusMessage
	p.repo.Update(ctx, p.db, project)

	p.stdLog.WarningFunction(fmt.Sprintf("Project %v scan failed", project.Key))
}

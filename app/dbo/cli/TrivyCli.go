package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sca-integrator/app/helper"
	"sca-integrator/app/shareVar"
	"strings"
)

type TrivyCli struct {
	command        string
	urlTrivyServer string
	trivyTimeout   string
	stdLog         *helper.StandartLog
}

func NewTrivyCli() *TrivyCli {
	return &TrivyCli{
		urlTrivyServer: os.Getenv("TRIVY_SERVER_URL"),
		trivyTimeout:   os.Getenv("TRIVY_TIMEOUT"),
		stdLog:         helper.NewStandardLog(shareVar.TrivyCli, shareVar.Repository),
	}
}

// trivy --timeout --server --ignore-unfixed --security-check vuln --skip-db-update --format json --output /_scanned-project-files/scan_repo_file_test_project_repository-.json '/_project-repository/workspace/test_project-sca'
// /_project-repository-file/test_project/
func (t *TrivyCli) Init() *TrivyCli {
	t.command = "trivy --timeout " + t.trivyTimeout + " --server " + t.urlTrivyServer
	return t
}

func (t *TrivyCli) AddFileSystemArgument(commandSkip string) *TrivyCli {
	if commandSkip != "" {
		t.command += " fs " + commandSkip
	} else {
		t.command += " fs"
	}
	return t
}

func (t *TrivyCli) AddIgnoreUnfixedArgument() *TrivyCli {
	t.command += " --ignore-unfixed"
	return t
}

func (t *TrivyCli) AddSeverityArgument(value string) *TrivyCli {
	if value != "" {
		t.command += " --severity " + value
	}
	return t
}

func (t *TrivyCli) AddSecurityCheckArgument(value string) *TrivyCli {
	t.command += " --security-check " + value
	return t
}

func (t *TrivyCli) AddSkipDBUpdateArgument() *TrivyCli {
	t.command += " --skip-db-update"
	return t
}

func (t *TrivyCli) AddFormat(value string) *TrivyCli {
	t.command += " --format " + value
	return t
}

func (t *TrivyCli) AddOutput(value string) *TrivyCli {
	t.command += fmt.Sprintf(" --output '%s'", value)
	return t
}

func (t *TrivyCli) AddProjectPath(value string) *TrivyCli {
	t.command += fmt.Sprintf(" '%s'", value)
	return t
}

func (t *TrivyCli) Exec(repositoryFilePath string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", t.command)
	cmd.Dir = repositoryFilePath
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	return out, err
}

func (t *TrivyCli) DownloadTrivyDb() bool {
	t.stdLog.InfoFunction("Start Download Vulnerability Database")

	successDownload := false
	cmd := exec.Command("bash", "-c", "trivy image --download-db-only")
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	if err != nil {
		stringOutput := fmt.Sprintf("%q\n", out)
		out, err = t.checkErrorIsTimeout(stringOutput, cmd, out, err)

		if err != nil {
			t.stdLog.WarningFunction("Download Vulnerability DB Failed")
		} else {
			//output := ""
			//if len(out) > 1 {
			//	output = "success downloading trivy db"
			//} else {
			//	output = "nothing to update, already up to date"
			//}
			t.stdLog.InfoFunction("End Download Vulnerability Database")
			successDownload = true
		}
	} else {
		//output := ""
		//if len(out) > 1 {
		//	output = "success downloading trivy db"
		//} else {
		//	output = "nothing to update, already up to date"
		//}
		t.stdLog.InfoFunction("End Download Vulnerability Database")
		successDownload = true
	}
	return successDownload
}

func (t *TrivyCli) checkErrorIsTimeout(errOutput string, cmd *exec.Cmd, out []byte, err error) ([]byte, error) {
	if strings.Contains(errOutput, "i/o timeout") {
		doAgain := true
		maxLoop := 5
		loop := 0

		for doAgain {
			loop++
			if loop > maxLoop {
				doAgain = false
			}
			out, err = cmd.CombinedOutput()
			if err != nil {
				errOutput = fmt.Sprintf("%q\n", out)
				if strings.Contains(errOutput, "i/o timeout") {
					doAgain = true
				} else {
					doAgain = false
				}
			} else {
				doAgain = false
			}
		}

		return out, err
	}
	return out, err
}

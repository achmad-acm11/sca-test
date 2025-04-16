package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"sca-integrator/app/dbo/cli"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/dbo/repository"
	"sca-integrator/app/dbo/repository/project"
	"sca-integrator/app/dto/request"
	"sca-integrator/app/dto/response"
	"sca-integrator/app/exception"
	"sca-integrator/app/helper"
	"sca-integrator/app/shareVar"
)

type ProjectService interface {
	GetAllData(ctx *gin.Context) []response.ProjectResponse
	GetDetailByIdData(ctx *gin.Context, id int) response.ProjectResponse
	ScanProjectData(ctx *gin.Context, request request.ProjectScanRequest)
	CreateProjectData(ctx *gin.Context, request request.CreateProjectRequest) response.ProjectResponse
	UpdateProjectData(ctx *gin.Context, request request.UpdateProjectRequest, id int) response.ProjectResponse
	DeleteProjectData(ctx *gin.Context, id int)
}

type ProjectServiceImpl struct {
	repo          repository.ProjectRepository
	repoOption    project.FilterOptionRepository
	repoExclusion project.ExclusionRepository
	repoResult    repository.ResultRepository
	trivyCli      *cli.TrivyCli
	validator     *validator.Validate
	db            *gorm.DB
	stdLog        *helper.StandartLog
}

func NewProjectService(repo repository.ProjectRepository,
	repoOption project.FilterOptionRepository,
	repoExclusion project.ExclusionRepository,
	repoResult repository.ResultRepository,
	trivyCli *cli.TrivyCli,
	validate *validator.Validate,
	db *gorm.DB) *ProjectServiceImpl {

	return &ProjectServiceImpl{
		repo:          repo,
		repoOption:    repoOption,
		repoExclusion: repoExclusion,
		repoResult:    repoResult,
		trivyCli:      trivyCli,
		validator:     validate,
		db:            db,
		stdLog:        helper.NewStandardLog(shareVar.Project, shareVar.Service),
	}
}

func (p ProjectServiceImpl) GetAllData(ctx *gin.Context) []response.ProjectResponse {
	projects := p.repo.GetAll(ctx, p.db)

	return response.NewProjectResponseBuilder().List(projects).ListResult()
}

func (p ProjectServiceImpl) GetDetailByIdData(ctx *gin.Context, id int) response.ProjectResponse {
	project := p.getProjectWithException(ctx, id)

	return response.NewProjectResponseBuilder().Default(project).Result()
}

func (p ProjectServiceImpl) CreateProjectData(ctx *gin.Context, request request.CreateProjectRequest) response.ProjectResponse {
	p.stdLog.NameFunc = "CreateProjectData"
	p.stdLog.StartFunction(request)

	err := p.validator.Struct(request)
	helper.ErrorHandlerValidator(err)

	tx := p.db.Begin()
	defer helper.CommitOrRollback(tx)

	project := p.repo.Create(ctx, tx, entity.Project{
		Key:         helper.ToSnakeCase(request.Name),
		Name:        request.Name,
		Description: request.Description,
		RepoType:    request.Repo_type,
		Url:         request.Url,
		BranchName:  request.Branch_name,
		Visibility:  request.Visibility,
	})

	p.stdLog.NameFunc = "CreateProjectData"
	p.stdLog.EndFunction(project)

	return response.NewProjectResponseBuilder().Default(project).Result()
}

func (p ProjectServiceImpl) UpdateProjectData(ctx *gin.Context, request request.UpdateProjectRequest, id int) response.ProjectResponse {
	err := p.validator.Struct(request)
	helper.ErrorHandlerValidator(err)

	tx := p.db.Begin()
	defer helper.CommitOrRollback(tx)

	project := p.getProjectWithException(ctx, id)

	project.Name = request.Name
	project.Description = request.Description
	project.RepoType = request.Repo_type
	project.Url = request.Url
	project.BranchName = request.Branch_name
	project.Visibility = request.Visibility

	projectNew := p.repo.Update(ctx, tx, project)

	return response.NewProjectResponseBuilder().Default(projectNew).Result()
}

func (p ProjectServiceImpl) DeleteProjectData(ctx *gin.Context, id int) {
	project := p.getProjectWithException(ctx, id)

	tx := p.db.Begin()
	defer helper.CommitOrRollback(tx)
	p.repo.DeleteOne(ctx, tx, project)
}

func (p ProjectServiceImpl) ScanProjectData(ctx *gin.Context, request request.ProjectScanRequest) {
	p.stdLog.NameFunc = "ScanProjectData"
	p.stdLog.StartFunction(request)

	err := p.validator.Struct(request)
	helper.ErrorHandlerValidator(err)

	project := p.getProjectWithException(ctx, request.ProjectId)
	p.stdLog.InfoFunction(project)

	p.checkStatusProject(project)

	project.StatusScan = 1
	p.repo.Update(ctx, p.db, project)

	switch request.ScanType {
	case "repository":
		go p.scanningRepository(ctx, project, request.Stage)
	case "image":
		panic(exception.NewNotImplementedError(errors.New(shareVar.NOT_IMPLEMENTED).Error()))
	default:
		panic(exception.NewNotImplementedError(errors.New(shareVar.NOT_IMPLEMENTED).Error()))
	}

	p.stdLog.NameFunc = "ScanProjectData"
	p.stdLog.EndFunction(project)
}

func (p ProjectServiceImpl) getProjectWithException(ctx *gin.Context, id int) entity.Project {
	project := p.repo.GetOneById(ctx, p.db, id)
	if project.Id == 0 {
		panic(exception.NewNotFoundError(errors.New(shareVar.PROJECT_NOT_FOUND).Error()))
	}

	return project
}

func (p ProjectServiceImpl) getOptionProjectWithException(ctx *gin.Context, projectId int) entity.ProjectFilterOption {
	option := p.repoOption.GetOneByProjectId(ctx, p.db, projectId)
	if option.Id == 0 {
		panic(exception.NewNotFoundError(errors.New(shareVar.FILTER_OPTION_NOT_FOUND).Error()))
	}

	return option
}

func (p ProjectServiceImpl) checkStatusProject(project entity.Project) {
	if project.StatusScan == 1 {
		panic(exception.NewConflictError(errors.New(shareVar.PROJECT_ON_SCANNING).Error()))
	}
}

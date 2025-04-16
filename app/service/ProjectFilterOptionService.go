package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/dbo/repository"
	"sca-integrator/app/dbo/repository/project"
	"sca-integrator/app/dto/request"
	"sca-integrator/app/dto/response"
	"sca-integrator/app/exception"
	"sca-integrator/app/helper"
	"sca-integrator/app/shareVar"
)

type ProjectFilterOptionService interface {
	GetAllData(ctx *gin.Context) []response.ProjectFilterOptionResponse
	GetDetailByIdData(ctx *gin.Context, id int) response.ProjectFilterOptionResponse
	CreateOptionData(ctx *gin.Context, request request.CreateOptionRequest) response.ProjectFilterOptionResponse
	UpdateOptionData(ctx *gin.Context, request request.UpdateOptionRequest, id int) response.ProjectFilterOptionResponse
	DeleteOptionData(ctx *gin.Context, id int)
}

type ProjectFilterOptionServiceImpl struct {
	repo        project.FilterOptionRepository
	repoProject repository.ProjectRepository
	validator   *validator.Validate
	db          *gorm.DB
}

func NewProjectFilterOptionService(repo project.FilterOptionRepository, repoProject repository.ProjectRepository, validate *validator.Validate, db *gorm.DB) *ProjectFilterOptionServiceImpl {
	return &ProjectFilterOptionServiceImpl{
		repo:        repo,
		repoProject: repoProject,
		validator:   validate,
		db:          db,
	}
}

func (p ProjectFilterOptionServiceImpl) GetAllData(ctx *gin.Context) []response.ProjectFilterOptionResponse {
	options := p.repo.GetAll(ctx, p.db)

	return response.NewProjectFilterOptionResponseBuilder().List(options).ListResult()
}

func (p ProjectFilterOptionServiceImpl) GetDetailByIdData(ctx *gin.Context, id int) response.ProjectFilterOptionResponse {
	option := p.GetOptionWithException(ctx, id)

	return response.NewProjectFilterOptionResponseBuilder().Default(option).Result()
}

func (p ProjectFilterOptionServiceImpl) CreateOptionData(ctx *gin.Context, request request.CreateOptionRequest) response.ProjectFilterOptionResponse {
	err := p.validator.Struct(request)
	helper.ErrorHandlerValidator(err)

	p.GetProjectWithException(ctx, request.Project_id)

	tx := p.db.Begin()
	defer helper.CommitOrRollback(tx)

	option := p.repo.Create(ctx, tx, entity.ProjectFilterOption{
		ProjectId:  request.Project_id,
		FilterType: request.Filter_type,
		Value:      request.Value,
	})

	return response.NewProjectFilterOptionResponseBuilder().Default(option).Result()
}

func (p ProjectFilterOptionServiceImpl) UpdateOptionData(ctx *gin.Context, request request.UpdateOptionRequest, id int) response.ProjectFilterOptionResponse {
	err := p.validator.Struct(request)
	helper.ErrorHandlerValidator(err)

	p.GetProjectWithException(ctx, request.Project_id)
	option := p.GetOptionWithException(ctx, id)

	tx := p.db.Begin()
	defer helper.CommitOrRollback(tx)

	option.ProjectId = request.Project_id
	option.FilterType = request.Filter_type
	option.Value = request.Value

	optionNew := p.repo.Update(ctx, tx, option)

	return response.NewProjectFilterOptionResponseBuilder().Default(optionNew).Result()
}

func (p ProjectFilterOptionServiceImpl) DeleteOptionData(ctx *gin.Context, id int) {
	option := p.GetOptionWithException(ctx, id)

	tx := p.db.Begin()
	defer helper.CommitOrRollback(tx)
	p.repo.DeleteOne(ctx, tx, option)
}

func (p ProjectFilterOptionServiceImpl) GetProjectWithException(ctx *gin.Context, projectId int) entity.Project {
	project := p.repoProject.GetOneById(ctx, p.db, projectId)
	if project.Id == 0 {
		panic(exception.NewNotFoundError(errors.New(shareVar.PROJECT_NOT_FOUND).Error()))
	}

	return project
}

func (p ProjectFilterOptionServiceImpl) GetOptionWithException(ctx *gin.Context, id int) entity.ProjectFilterOption {
	option := p.repo.GetOneById(ctx, p.db, id)
	if option.Id == 0 {
		panic(exception.NewNotFoundError(errors.New(shareVar.FILTER_OPTION_NOT_FOUND).Error()))
	}

	return option
}

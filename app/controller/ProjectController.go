package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sca-integrator/app/dto/request"
	"sca-integrator/app/helper"
	"sca-integrator/app/service"
	"sca-integrator/app/shareVar"
	"strconv"
)

type ProjectController struct {
	service service.ProjectService
}

func NewProjectController(service service.ProjectService) *ProjectController {
	return &ProjectController{
		service: service,
	}
}

func (p *ProjectController) GetAllHandler(ctx *gin.Context) {
	responses := p.service.GetAllData(ctx)

	ctx.JSON(http.StatusOK, responses)
}

func (p *ProjectController) GetDetailByIdHandler(ctx *gin.Context) {
	projectId, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandlerValidator(err)

	response := p.service.GetDetailByIdData(ctx, projectId)

	ctx.JSON(http.StatusOK, response)
}

func (p *ProjectController) ScanProjectHandler(ctx *gin.Context) {
	var request request.ProjectScanRequest

	err := ctx.ShouldBindJSON(&request)
	helper.ErrorHandler(err)

	p.service.ScanProjectData(ctx, request)

	ctx.JSON(http.StatusOK, shareVar.PROJECT_SCAN_STARTED)
}

func (p *ProjectController) CreateProjectHandler(ctx *gin.Context) {
	var request request.CreateProjectRequest
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorHandler(err)

	response := p.service.CreateProjectData(ctx, request)

	ctx.JSON(http.StatusOK, response)
}

func (p *ProjectController) UpdateProjectHandler(ctx *gin.Context) {
	projectId, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	var request request.UpdateProjectRequest
	err = ctx.ShouldBindJSON(&request)
	helper.ErrorHandler(err)

	response := p.service.UpdateProjectData(ctx, request, projectId)

	ctx.JSON(http.StatusOK, response)
}

func (p *ProjectController) DeleteProjectHandler(ctx *gin.Context) {
	projectId, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	p.service.DeleteProjectData(ctx, projectId)

	ctx.JSON(http.StatusOK, shareVar.PROJECT_DELETED)
}

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

type ProjectFilterOptionController struct {
	service service.ProjectFilterOptionService
}

func NewProjectFilterOptionController(service service.ProjectFilterOptionService) *ProjectFilterOptionController {
	return &ProjectFilterOptionController{
		service: service,
	}
}

func (p *ProjectFilterOptionController) GetAllHandler(ctx *gin.Context) {
	responses := p.service.GetAllData(ctx)

	ctx.JSON(http.StatusOK, responses)
}

func (p *ProjectFilterOptionController) GetDetailByIdHandler(ctx *gin.Context) {
	optionId, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandlerValidator(err)

	response := p.service.GetDetailByIdData(ctx, optionId)

	ctx.JSON(http.StatusOK, response)
}

func (p *ProjectFilterOptionController) CreateOptionHandler(ctx *gin.Context) {
	var request request.CreateOptionRequest
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorHandler(err)

	response := p.service.CreateOptionData(ctx, request)

	ctx.JSON(http.StatusOK, response)
}

func (p *ProjectFilterOptionController) UpdateOptionHandler(ctx *gin.Context) {
	optionId, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	var request request.UpdateOptionRequest
	err = ctx.ShouldBindJSON(&request)
	helper.ErrorHandler(err)

	response := p.service.UpdateOptionData(ctx, request, optionId)

	ctx.JSON(http.StatusOK, response)
}

func (p *ProjectFilterOptionController) DeleteOptionHandler(ctx *gin.Context) {
	optionId, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	p.service.DeleteOptionData(ctx, optionId)

	ctx.JSON(http.StatusOK, shareVar.FILTER_OPTION_DELETED)
}

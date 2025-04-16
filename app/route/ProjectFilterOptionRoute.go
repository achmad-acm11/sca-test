package route

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"sca-integrator/app/controller"
	"sca-integrator/app/dbo/repository"
	"sca-integrator/app/dbo/repository/project"
	"sca-integrator/app/service"
)

func OptionRoute(router *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) *gin.RouterGroup {
	optionRepository := project.NewFilterOptionRepository()
	projectRepository := repository.NewProjectRepository()

	optionService := service.NewProjectFilterOptionService(optionRepository, projectRepository, validate, db)

	optionController := controller.NewProjectFilterOptionController(optionService)

	router.GET("filter-options", optionController.GetAllHandler)
	router.GET("filter-option/:id", optionController.GetDetailByIdHandler)
	router.POST("filter-option", optionController.CreateOptionHandler)
	router.PUT("filter-option/:id", optionController.UpdateOptionHandler)
	router.DELETE("filter-option/:id", optionController.DeleteOptionHandler)

	return router
}

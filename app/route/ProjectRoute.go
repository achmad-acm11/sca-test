package route

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"sca-integrator/app/controller"
	"sca-integrator/app/dbo/cli"
	"sca-integrator/app/dbo/repository"
	"sca-integrator/app/dbo/repository/project"
	"sca-integrator/app/service"
)

func ProjectRoute(router *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) *gin.RouterGroup {
	projectRepository := repository.NewProjectRepository()
	exclusionRepository := project.NewExclusionRepository()
	optionRepository := project.NewFilterOptionRepository()
	resultRepository := repository.NewResultRepository()
	trivyCli := cli.NewTrivyCli()

	projectService := service.NewProjectService(projectRepository,
		optionRepository,
		exclusionRepository,
		resultRepository,
		trivyCli,
		validate,
		db)

	projectController := controller.NewProjectController(projectService)

	router.GET("projects", projectController.GetAllHandler)
	router.GET("project/:id", projectController.GetDetailByIdHandler)
	router.POST("project", projectController.CreateProjectHandler)
	router.POST("project/scan", projectController.ScanProjectHandler)
	router.PUT("project/:id", projectController.UpdateProjectHandler)
	router.DELETE("project/:id", projectController.DeleteProjectHandler)

	return router
}

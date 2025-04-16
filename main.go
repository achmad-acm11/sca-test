package main

import (
	"github.com/gin-gonic/gin"
	validator2 "github.com/go-playground/validator/v10"
	cors "github.com/itsjamie/gin-cors"
	"github.com/joho/godotenv"
	"os"
	"sca-integrator/app/config"
	"sca-integrator/app/middleware"
	"sca-integrator/app/migration"
	"sca-integrator/app/route"
	"sca-integrator/app/shareVar"
	"time"
)

func main() {
	if os.Getenv("APP_ENV") == "" {
		errEnv := godotenv.Load(".env")

		if errEnv != nil {
			panic("Failed load file .env")
		}
	}

	app_port := os.Getenv("APP_PORT")
	db := config.ConfigDB()
	migration.DoMigration(db)
	validator := validator2.New()
	shareVar.Logger = config.ConfigLog()

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		ValidateHeaders: false,
	}))

	api := router.Group("api/v1/")

	api = route.ProjectRoute(api, db, validator)
	api = route.OptionRoute(api, db, validator)

	err := router.Run(":" + app_port)
	if err != nil {
		panic("Failed run application")
	}
}

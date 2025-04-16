package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"sca-integrator/app/helper"
)

func ConfigDB() *gorm.DB {
	if os.Getenv("APP_ENV") == "" {
		errEnv := godotenv.Load(".env")
		helper.ErrorHandler(errEnv)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	config := fmt.Sprintf(
		//"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, db_name)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	helper.ErrorHandler(err)

	return db
}

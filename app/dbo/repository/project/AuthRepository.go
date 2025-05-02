package project

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/helper"
)

type AuthRepository interface {
	Create(ctx *gin.Context, db *gorm.DB, projectAuth entity.ProjectAuth) entity.ProjectAuth
	Update(ctx *gin.Context, db *gorm.DB, projectAuth entity.ProjectAuth) entity.ProjectAuth
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.ProjectAuth
	DeleteOne(ctx *gin.Context, db *gorm.DB, projectAuth entity.ProjectAuth)
}
type AuthRepositoryImpl struct {
}

func NewAuthRepository() *AuthRepositoryImpl {
	return &AuthRepositoryImpl{}
}

func (r *AuthRepositoryImpl) Create(ctx *gin.Context, db *gorm.DB, projectAuth entity.ProjectAuth) entity.ProjectAuth {
	err := db.WithContext(ctx).Create(&projectAuth).Error
	helper.ErrorHandler(err)

	return projectAuth
}

func (r *AuthRepositoryImpl) Update(ctx *gin.Context, db *gorm.DB, projectAuth entity.ProjectAuth) entity.ProjectAuth {
	err := db.WithContext(ctx).Model(&projectAuth).Updates(projectAuth).Error
	helper.ErrorHandler(err)

	return projectAuth
}

func (r *AuthRepositoryImpl) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.ProjectAuth {
	projectAuth := entity.ProjectAuth{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&projectAuth).Error
	helper.ErrorHandler(err)

	return projectAuth
}

func (r *AuthRepositoryImpl) DeleteOne(ctx *gin.Context, db *gorm.DB, projectAuth entity.ProjectAuth) {
	err := db.WithContext(ctx).Delete(&projectAuth).Error
	helper.ErrorHandler(err)
}

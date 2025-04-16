package repository

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/helper"
	"time"
)

type ProjectRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entity.Project
	Create(ctx *gin.Context, db *gorm.DB, project entity.Project) entity.Project
	Update(ctx *gin.Context, db *gorm.DB, project entity.Project) entity.Project
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.Project
	DeleteOne(ctx *gin.Context, db *gorm.DB, project entity.Project)
}

type ProjectRepositoryImpl struct {
}

func NewProjectRepository() *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{}
}

func (p ProjectRepositoryImpl) GetAll(ctx *gin.Context, db *gorm.DB) []entity.Project {
	projects := []entity.Project{}

	tx := db.WithContext(ctx)
	err := tx.Find(&projects).Error
	helper.ErrorHandler(err)

	return projects
}

func (p ProjectRepositoryImpl) Create(ctx *gin.Context, db *gorm.DB, project entity.Project) entity.Project {
	err := db.WithContext(ctx).Create(&project).Error
	helper.ErrorHandler(err)

	return project
}

func (p ProjectRepositoryImpl) Update(ctx *gin.Context, db *gorm.DB, project entity.Project) entity.Project {
	project.UpdatedAt = time.Now()
	
	err := db.WithContext(ctx).Model(&project).Updates(project).Error
	helper.ErrorHandler(err)

	return project
}

func (p ProjectRepositoryImpl) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.Project {
	project := entity.Project{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&project).Error
	helper.ErrorHandler(err)

	return project
}

func (p ProjectRepositoryImpl) DeleteOne(ctx *gin.Context, db *gorm.DB, project entity.Project) {
	err := db.WithContext(ctx).Delete(&project).Error
	helper.ErrorHandler(err)
}

package project

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/helper"
	"time"
)

type FilterOptionRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entity.ProjectFilterOption
	GetAllByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) []entity.ProjectFilterOption
	Create(ctx *gin.Context, db *gorm.DB, option entity.ProjectFilterOption) entity.ProjectFilterOption
	Update(ctx *gin.Context, db *gorm.DB, option entity.ProjectFilterOption) entity.ProjectFilterOption
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.ProjectFilterOption
	GetOneByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) entity.ProjectFilterOption
	GetAllByProjectIdAndFilterType(ctx *gin.Context, db *gorm.DB, projectId int, filterType string) []entity.ProjectFilterOption
	DeleteOne(ctx *gin.Context, db *gorm.DB, option entity.ProjectFilterOption)
}

type FilterOptionRepositoryImpl struct {
}

func NewFilterOptionRepository() *FilterOptionRepositoryImpl {
	return &FilterOptionRepositoryImpl{}
}

func (p FilterOptionRepositoryImpl) GetAll(ctx *gin.Context, db *gorm.DB) []entity.ProjectFilterOption {
	options := []entity.ProjectFilterOption{}

	tx := db.WithContext(ctx)
	err := tx.Find(&options).Error

	helper.ErrorHandler(err)

	return options
}

func (p FilterOptionRepositoryImpl) GetAllByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) []entity.ProjectFilterOption {
	options := []entity.ProjectFilterOption{}

	tx := db.WithContext(ctx)
	err := tx.Where("project_id = ?", projectId).Find(&options).Error

	helper.ErrorHandler(err)

	return options
}

func (p FilterOptionRepositoryImpl) Create(ctx *gin.Context, db *gorm.DB, option entity.ProjectFilterOption) entity.ProjectFilterOption {
	err := db.WithContext(ctx).Create(&option).Error

	helper.ErrorHandler(err)

	return option
}

func (p FilterOptionRepositoryImpl) Update(ctx *gin.Context, db *gorm.DB, option entity.ProjectFilterOption) entity.ProjectFilterOption {
	option.UpdatedAt = time.Now()

	err := db.WithContext(ctx).Model(&option).Updates(option).Error

	helper.ErrorHandler(err)

	return option
}

func (p FilterOptionRepositoryImpl) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.ProjectFilterOption {
	option := entity.ProjectFilterOption{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&option).Error
	helper.ErrorHandler(err)

	return option
}

func (p FilterOptionRepositoryImpl) DeleteOne(ctx *gin.Context, db *gorm.DB, option entity.ProjectFilterOption) {
	err := db.WithContext(ctx).Delete(&option).Error
	helper.ErrorHandler(err)
}

func (p FilterOptionRepositoryImpl) GetOneByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) entity.ProjectFilterOption {
	option := entity.ProjectFilterOption{}

	err := db.WithContext(ctx).Where("project_id = ?", projectId).Find(&option).Error
	helper.ErrorHandler(err)

	return option
}

func (p FilterOptionRepositoryImpl) GetAllByProjectIdAndFilterType(ctx *gin.Context, db *gorm.DB, projectId int, filterType string) []entity.ProjectFilterOption {
	options := []entity.ProjectFilterOption{}

	err := db.WithContext(ctx).Where("project_id = ?", projectId).Where("filter_type = ?", filterType).Find(&options).Error
	helper.ErrorHandler(err)

	return options
}

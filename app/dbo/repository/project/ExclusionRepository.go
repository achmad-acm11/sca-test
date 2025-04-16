package project

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/helper"
)

type ExclusionRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entity.ProjectExclusion
	Create(ctx *gin.Context, db *gorm.DB, exclusion entity.ProjectExclusion) entity.ProjectExclusion
	Update(ctx *gin.Context, db *gorm.DB, exclusion entity.ProjectExclusion) entity.ProjectExclusion
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.ProjectExclusion
	GetOneByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) entity.ProjectExclusion
	GetAllByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) []entity.ProjectExclusion
	DeleteOne(ctx *gin.Context, db *gorm.DB, exclusion entity.ProjectExclusion)
}

type ExclusionRepositoryImpl struct {
}

func NewExclusionRepository() *ExclusionRepositoryImpl {
	return &ExclusionRepositoryImpl{}
}

func (e ExclusionRepositoryImpl) GetAll(ctx *gin.Context, db *gorm.DB) []entity.ProjectExclusion {
	exclusions := []entity.ProjectExclusion{}

	tx := db.WithContext(ctx)
	err := tx.Find(&exclusions).Error
	helper.ErrorHandler(err)

	return exclusions
}

func (e ExclusionRepositoryImpl) Create(ctx *gin.Context, db *gorm.DB, exclusion entity.ProjectExclusion) entity.ProjectExclusion {
	err := db.WithContext(ctx).Create(&exclusion).Error
	helper.ErrorHandler(err)

	return exclusion
}

func (e ExclusionRepositoryImpl) Update(ctx *gin.Context, db *gorm.DB, exclusion entity.ProjectExclusion) entity.ProjectExclusion {
	err := db.WithContext(ctx).Model(&exclusion).Updates(exclusion).Error
	helper.ErrorHandler(err)

	return exclusion
}

func (e ExclusionRepositoryImpl) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.ProjectExclusion {
	exclusion := entity.ProjectExclusion{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&exclusion).Error
	helper.ErrorHandler(err)

	return exclusion
}

func (e ExclusionRepositoryImpl) GetOneByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) entity.ProjectExclusion {
	exclusion := entity.ProjectExclusion{}

	err := db.WithContext(ctx).Where("project_id = ?", projectId).Find(&exclusion).Error
	helper.ErrorHandler(err)

	return exclusion
}

func (e ExclusionRepositoryImpl) GetAllByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) []entity.ProjectExclusion {
	exclusions := []entity.ProjectExclusion{}

	err := db.WithContext(ctx).Where("project_id = ?", projectId).Find(&exclusions).Error
	helper.ErrorHandler(err)

	return exclusions
}

func (e ExclusionRepositoryImpl) DeleteOne(ctx *gin.Context, db *gorm.DB, exclusion entity.ProjectExclusion) {
	err := db.WithContext(ctx).Delete(&exclusion).Error
	helper.ErrorHandler(err)
}

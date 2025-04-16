package repository

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sca-integrator/app/dbo/entity"
	"sca-integrator/app/helper"
)

type ResultRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entity.Result
	Create(ctx *gin.Context, db *gorm.DB, result entity.Result) entity.Result
	Update(ctx *gin.Context, db *gorm.DB, result entity.Result) entity.Result
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.Result
	DeleteOne(ctx *gin.Context, db *gorm.DB, result entity.Result)
	GetLastByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) entity.Result
	GetOneResultByProjectIdAndRuleAndTargetFile(ctx *gin.Context, db *gorm.DB, projectId int, rule string, targetFile string) entity.Result
	GetAllByProjectIdAndScanVersion(ctx *gin.Context, db *gorm.DB, projectId int, scanVersion int) []entity.Result
}

type ResultRepositoryImpl struct {
}

func NewResultRepository() *ResultRepositoryImpl {
	return &ResultRepositoryImpl{}
}

func (r ResultRepositoryImpl) GetAll(ctx *gin.Context, db *gorm.DB) []entity.Result {
	results := []entity.Result{}

	tx := db.WithContext(ctx)
	err := tx.Find(&results).Error

	helper.ErrorHandler(err)

	return results
}

func (r ResultRepositoryImpl) Create(ctx *gin.Context, db *gorm.DB, result entity.Result) entity.Result {
	err := db.WithContext(ctx).Create(&result).Error

	helper.ErrorHandler(err)

	return result
}

func (r ResultRepositoryImpl) Update(ctx *gin.Context, db *gorm.DB, result entity.Result) entity.Result {
	err := db.WithContext(ctx).Model(&result).Updates(result).Error

	helper.ErrorHandler(err)

	return result
}

func (r ResultRepositoryImpl) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entity.Result {
	result := entity.Result{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&result).Error
	helper.ErrorHandler(err)

	return result
}

func (r ResultRepositoryImpl) DeleteOne(ctx *gin.Context, db *gorm.DB, result entity.Result) {
	err := db.WithContext(ctx).Delete(&result).Error
	helper.ErrorHandler(err)
}

func (r ResultRepositoryImpl) GetLastByProjectId(ctx *gin.Context, db *gorm.DB, projectId int) entity.Result {
	result := entity.Result{}

	err := db.WithContext(ctx).Where("project_id = ?", projectId).Find(&result).Error
	helper.ErrorHandler(err)

	return result
}

func (r ResultRepositoryImpl) GetOneResultByProjectIdAndRuleAndTargetFile(ctx *gin.Context, db *gorm.DB, projectId int, rule string, targetFile string) entity.Result {
	result := entity.Result{}

	err := db.WithContext(ctx).
		Where("project_id = ?", projectId).
		Where("rule = ?", rule).
		Where("target_file = ?", targetFile).
		Find(&result).Error
	helper.ErrorHandler(err)

	return result
}

func (r ResultRepositoryImpl) GetAllByProjectIdAndScanVersion(ctx *gin.Context, db *gorm.DB, projectId int, scanVersion int) []entity.Result {
	results := []entity.Result{}

	tx := db.WithContext(ctx)
	err := tx.Where("project_id = ?", projectId).
		Where("scan_version = ?", scanVersion).
		Find(&results).Error

	helper.ErrorHandler(err)

	return results
}

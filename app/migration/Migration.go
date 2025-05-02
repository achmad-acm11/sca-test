package migration

import (
	"gorm.io/gorm"
	"sca-integrator/app/dbo/entity"
)

func DoMigration(db *gorm.DB) {
	db.AutoMigrate(&entity.Project{})
	db.AutoMigrate(&entity.ProjectFilterOption{})
	db.AutoMigrate(&entity.ProjectExclusion{})
	db.AutoMigrate(&entity.Result{})
	db.AutoMigrate(&entity.ProjectAuth{})
}

package entity

import (
	"gorm.io/gorm"
	"time"
)

type ProjectExclusion struct {
	Id        int            `gorm:"column:id;type:int;primaryKey;autoIncrement;not null"`
	ProjectId int            `gorm:"column:project_id;type:int"`
	Path      string         `gorm:"column:path;type:varchar(255)"`
	Type      string         `gorm:"column:type;type:varchar(255)"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;->"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;default:null;->"`
}

func (ProjectExclusion) TableName() string {
	return "project_exclusions"
}

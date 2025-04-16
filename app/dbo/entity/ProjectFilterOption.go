package entity

import (
	"gorm.io/gorm"
	"time"
)

type ProjectFilterOption struct {
	Id          int            `gorm:"column:id;type:int;primaryKey;autoIncrement;not null"`
	ProjectId   int            `gorm:"column:project_id;type:int"`
	FilterType  string         `gorm:"column:filter_type;type:varchar(255)"`
	Value       string         `gorm:"column:value;type:varchar(255)"`
	ScanVersion int            `gorm:"column:scan_version;type:int"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;->"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;default:null;->"`
}

func (ProjectFilterOption) TableName() string {
	return "project_filter_options"
}

package entity

import (
	"gorm.io/gorm"
	"time"
)

type Project struct {
	Id                 int            `gorm:"column:id;type:int;primaryKey;autoIncrement;not null"`
	Key                string         `gorm:"column:key;type:varchar(255)"`
	Name               string         `gorm:"column:name;type:varchar(255)"`
	Description        string         `gorm:"column:description;type:varchar(255)"`
	RepoType           string         `gorm:"column:repo_type;type:varchar(255)"`
	Url                string         `gorm:"column:url;type:varchar(255)"`
	BranchName         string         `gorm:"column:branch_name;type:varchar(255)"`
	Visibility         string         `gorm:"column:visibility;type:varchar(255);default:PUBLIC"`
	StatusScan         int            `gorm:"column:status_scan;type:int;default:0"`
	CurrentScanVersion int            `gorm:"column:current_scan_version;type:int;default:0"`
	CreatedAt          time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;->"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;default:null;->"`
}

func (Project) TableName() string {
	return "projects"
}

package entity

import (
	"gorm.io/gorm"
	"time"
)

type ProjectAuth struct {
	Id        int            `gorm:"column:id;type:int;primaryKey;autoIncrement;not null"`
	ProjectId int            `gorm:"column:project_id;type:int;not null"`
	Username  string         `gorm:"column:username;type:varchar(64);not null"`
	Token     string         `gorm:"column:token;type:varchar(255);not null"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;->"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;default:null;->"`
}

func (p *ProjectAuth) TableName() string {
	return "project_auths"
}

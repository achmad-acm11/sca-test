package entity

import (
	"gorm.io/gorm"
	"time"
)

type Result struct {
	Id               int            `gorm:"column:id;type:int;primaryKey;autoIncrement;not null"`
	ProjectId        int            `gorm:"column:project_id;type:int"`
	Rule             string         `gorm:"column:rule;type:varchar(255)"`
	PrimaryUrl       string         `gorm:"column:primary_url;type:varchar(255)"`
	TargetFile       string         `gorm:"column:target_file;type:varchar(255)"`
	ScanType         string         `gorm:"column:scan_type;type:varchar(255)"`
	PackageName      string         `gorm:"column:package_name;type:varchar(255)"`
	Title            string         `gorm:"column:title;type:varchar(255)"`
	Description      string         `gorm:"column:description;type:varchar(255)"`
	Severity         string         `gorm:"column:severity;type:varchar(255)"`
	LastFoundAt      string         `gorm:"column:last_found_at;type:varchar(255)"`
	StatusResult     int            `gorm:"column:status_result;type:int"`
	PublishedDate    time.Time      `gorm:"column:published_date;type:timestamp"`
	LastModifiedDate time.Time      `gorm:"column:last_modified_date;type:timestamp"`
	CvssSource       string         `gorm:"column:cvss_source;type:varchar(255)"`
	CvssV2           string         `gorm:"column:cvss_v2;type:varchar(255)"`
	CvssV3           string         `gorm:"column:cvss_v3;type:varchar(255)"`
	InstalledVersion string         `gorm:"column:installed_version;type:varchar(255)"`
	FixedVersion     string         `gorm:"column:fixed_version;type:varchar(255)"`
	References       string         `gorm:"column:references;type:varchar(255)"`
	PackagesType     string         `gorm:"column:packages_type;type:varchar(255)"`
	ScanVersion      int            `gorm:"column:scan_version;type:int"`
	CreatedAt        time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;->"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;default:null;->"`
}

func (Result) TableName() string {
	return "results"
}

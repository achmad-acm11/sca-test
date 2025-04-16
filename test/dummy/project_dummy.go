package dummy

import (
	"database/sql/driver"
	"sca-integrator/app/dbo/entity"
)

var ProjectCols = []string{
	"id", "key", "name", "description", "repo_type", "url", "branch_name", "visibility", "status_scan",
	"current_scan_version"}

var ProjectDummyList = []entity.Project{
	entity.Project{
		Key:                "ini_project",
		Name:               "Ini Project",
		Description:        "Ini Desc",
		RepoType:           "github",
		Url:                "https://github.com/ini_project",
		BranchName:         "main",
		Visibility:         "PUBLIC",
		StatusScan:         0,
		CurrentScanVersion: 0,
	},
}

func MappingProjectStore(item entity.Project, numId int) []driver.Value {
	var values []driver.Value

	values = append(values, numId)
	values = append(values, item.Key)
	values = append(values, item.Name)
	values = append(values, item.Description)
	values = append(values, item.RepoType)
	values = append(values, item.Url)
	values = append(values, item.BranchName)
	values = append(values, item.Visibility)
	values = append(values, item.StatusScan)
	values = append(values, item.CurrentScanVersion)

	return values
}

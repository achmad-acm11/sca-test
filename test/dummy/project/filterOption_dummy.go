package project

import (
	"database/sql/driver"
	"sca-integrator/app/dbo/entity"
)

var FilterOptionCols = []string{
	"id", "project_id", "filter_type", "value", "scan_version"}

var FilterOptionDummyList = []entity.ProjectFilterOption{
	entity.ProjectFilterOption{
		ProjectId:   1,
		FilterType:  "Ini Filter",
		Value:       "ini filter",
		ScanVersion: 1,
	},
}

func MappingFilterOptionStore(item entity.ProjectFilterOption, numId int) []driver.Value {
	var values []driver.Value

	values = append(values, numId)
	values = append(values, item.ProjectId)
	values = append(values, item.FilterType)
	values = append(values, item.Value)
	values = append(values, item.ScanVersion)

	return values
}

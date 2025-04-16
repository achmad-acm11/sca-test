package response

import "sca-integrator/app/dbo/entity"

type ProjectFilterOptionResponse struct {
	Id         int    `json:"id"`
	ProjectId  int    `json:"project_id"`
	FilterType string `json:"filter_type"`
	Value      string `json:"value"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type ProjectFilterOptionResponseBuilder struct {
	singleData ProjectFilterOptionResponse
	listData   []ProjectFilterOptionResponse
}

func NewProjectFilterOptionResponseBuilder() *ProjectFilterOptionResponseBuilder {
	return &ProjectFilterOptionResponseBuilder{}
}

func (b *ProjectFilterOptionResponseBuilder) Default(option entity.ProjectFilterOption) *ProjectFilterOptionResponseBuilder {
	b.singleData = mapOption(option)
	return b
}

func (b *ProjectFilterOptionResponseBuilder) List(list []entity.ProjectFilterOption) *ProjectFilterOptionResponseBuilder {
	b.listData = mapListOption(list)
	return b
}

func (b *ProjectFilterOptionResponseBuilder) Result() ProjectFilterOptionResponse {
	return b.singleData
}

func (b *ProjectFilterOptionResponseBuilder) ListResult() []ProjectFilterOptionResponse {
	return b.listData
}

func mapOption(option entity.ProjectFilterOption) ProjectFilterOptionResponse {
	response := ProjectFilterOptionResponse{
		Id:         option.Id,
		ProjectId:  option.ProjectId,
		FilterType: option.FilterType,
		Value:      option.Value,
		CreatedAt:  option.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  option.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return response
}

func mapListOption(list []entity.ProjectFilterOption) []ProjectFilterOptionResponse {
	responses := []ProjectFilterOptionResponse{}
	for _, option := range list {
		responses = append(responses, mapOption(option))
	}
	return responses
}

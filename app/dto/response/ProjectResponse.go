package response

import (
	"sca-integrator/app/dbo/entity"
)

type ProjectResponse struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RepoType    string `json:"repo_type"`
	Url         string `json:"url"`
	BranchName  string `json:"branch_name"`
	Visibility  string `json:"visibility"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ProjectResponseBuilder struct {
	singleData ProjectResponse
	listData   []ProjectResponse
}

func NewProjectResponseBuilder() *ProjectResponseBuilder {
	return &ProjectResponseBuilder{}
}

func (builder *ProjectResponseBuilder) Default(project entity.Project) *ProjectResponseBuilder {
	builder.singleData = mapProject(project)
	return builder
}

func (builder *ProjectResponseBuilder) List(projects []entity.Project) *ProjectResponseBuilder {
	builder.listData = mapListProject(projects)
	return builder
}

func (builder *ProjectResponseBuilder) Result() ProjectResponse {
	return builder.singleData
}

func (builder *ProjectResponseBuilder) ListResult() []ProjectResponse {
	return builder.listData
}

func mapProject(project entity.Project) ProjectResponse {
	response := ProjectResponse{
		Id:          project.Id,
		Key:         project.Key,
		Name:        project.Name,
		Description: project.Description,
		RepoType:    project.RepoType,
		Url:         project.Url,
		BranchName:  project.BranchName,
		Visibility:  project.Visibility,
		CreatedAt:   project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return response
}

func mapListProject(projectList []entity.Project) []ProjectResponse {
	responses := []ProjectResponse{}
	for _, project := range projectList {
		responses = append(responses, mapProject(project))
	}
	return responses
}

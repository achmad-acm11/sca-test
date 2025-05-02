package request

type CreateProjectRequest struct {
	Name        string `validate:"required" json:"name"`
	Description string `json:"description"`
	Repo_type   string `validate:"required" json:"repo_type"`
	Url         string `validate:"required" json:"url"`
	Branch_name string `validate:"required" json:"branch_name"`
	Visibility  string `validate:"required,oneof=PUBLIC PRIVATE'" json:"visibility"`
	Username    string `json:"username"`
	Token       string `json:"token"`
}

type UpdateProjectRequest struct {
	Name        string `validate:"required" json:"name"`
	Description string `json:"description"`
	Repo_type   string `validate:"required" json:"repo_type"`
	Url         string `validate:"required" json:"url"`
	Branch_name string `validate:"required" json:"branch_name"`
	Visibility  string `validate:"required,oneof=PUBLIC PRIVATE'" json:"visibility"`
}

type ProjectScanRequest struct {
	ProjectId int    `validate:"required,numeric" json:"project_id"`
	ScanType  string `validate:"required" json:"scan_type"`
	Stage     string `json:"stage"`
	//OrganizationId int32  `json:"organization_id"`
}

package request

type CreateOptionRequest struct {
	Project_id  int    `validate:"required" json:"project_id"`
	Filter_type string `validate:"required" json:"filter_type"`
	Value       string `validate:"required" json:"value"`
}

type UpdateOptionRequest struct {
	Project_id  int    `validate:"required" json:"project_id"`
	Filter_type string `validate:"required" json:"filter_type"`
	Value       string `validate:"required" json:"value"`
}

package request

type GitCloneRequest struct {
	RepoUrl    string
	BranchName string
	RepoName   string
	Visibility string
	Username   string
	Token      string
}

package service

import (
	"net/url"
	"sca-integrator/app/dbo/cli"
	"sca-integrator/app/dto/request"
	"sca-integrator/app/helper"
	"sca-integrator/app/shareVar"
)

type GitCliService interface {
	RunCloneRepo(request request.GitCloneRequest)
}

type GitCliServiceImpl struct {
	GitCli *cli.GitCli
	stdLog *helper.StandartLog
}

func NewGitCliService(gitCli *cli.GitCli) *GitCliServiceImpl {
	return &GitCliServiceImpl{
		GitCli: gitCli,
		stdLog: helper.NewStandardLog(shareVar.GitCli, shareVar.Service),
	}
}

func (g *GitCliServiceImpl) RunCloneRepo(request request.GitCloneRequest) {
	g.stdLog.NameFunc = "RunCloneRepo"
	g.stdLog.StartFunction(request)

	if request.Visibility == "PRIVATE" {
		request.RepoUrl = g.modifyUrlForPrivate(request.RepoUrl, request.Username, request.Token)
	}

	g.GitCli.InitClone(request.RepoUrl).
		AddParamBranch(request.BranchName).
		RepoDir(request.RepoName).
		Exec()

	g.stdLog.NameFunc = "RunCloneRepo"
	g.stdLog.EndFunction(nil)
}

func (g *GitCliServiceImpl) modifyUrlForPrivate(repoUrl string, username string, token string) string {
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		panic(err)
	}

	parsedURL.User = url.UserPassword(username, token)

	return parsedURL.String()
}

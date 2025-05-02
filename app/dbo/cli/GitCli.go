package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sca-integrator/app/helper"
	"sca-integrator/app/shareVar"
)

type GitCli struct {
	repoPath string
	command  string
	stdLog   *helper.StandartLog
}

func NewGitCli() *GitCli {
	var curDir, _ = os.Getwd()
	repoPath := fmt.Sprintf("%s/_project-repository/workspace/", curDir)
	if _, err := os.Stat(repoPath); !os.IsExist(err) {
		fmt.Printf("Create Repo %s \n", repoPath)
		err := os.MkdirAll(repoPath, 0755)
		fmt.Println(err)
	}

	return &GitCli{
		repoPath: repoPath,
		stdLog:   helper.NewStandardLog(shareVar.GitCli, shareVar.Repository),
	}
}

func (g *GitCli) InitClone(repoUrl string) *GitCli {
	g.command = "git clone " + repoUrl
	return g
}

func (g *GitCli) AddParamBranch(branch string) *GitCli {
	g.command += " -b " + branch
	return g
}

func (g *GitCli) RepoDir(repoName string) *GitCli {
	repoDir := g.repoPath + repoName + "-sca"
	if _, err := os.Stat(repoDir); !os.IsExist(err) {
		fmt.Printf("Create Repo Dir %s \n", repoDir)
		err := os.Mkdir(repoDir, 0755)
		fmt.Println(err)
	}

	g.command += " " + repoDir
	return g
}

func (g *GitCli) Exec() {
	g.stdLog.InfoFunction("Executing command: " + g.command)

	cmd := exec.Command("bash", "-c", g.command)

	if err := cmd.Run(); err != nil {
		g.stdLog.ErrorFunction("Error executing command: " + g.command)
	}
}

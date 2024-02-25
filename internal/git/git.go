package git

import (
	"errors"
	"os/exec"

	"github.com/go-git/go-git/v5"
)

var (
	Worktree     *git.Worktree
	Repo         *git.Repository
	GitInstalled bool
	GitBin       string
)

func Init() error {
	var err error
	GitBin, err = exec.LookPath("git")
	if err != nil || GitBin == "" {
		// git not found
		GitInstalled = false
	} else {
		// git found
		GitInstalled = true
	}

	opt := git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: true}
	Repo, err = git.PlainOpenWithOptions(".", &opt)
	if err != nil {
		return errors.Join(errors.New("failed to open git repository"), err)
	}

	Worktree, err = Repo.Worktree()
	if err != nil {
		return errors.Join(errors.New("failed to get worktree"), err)
	}

	return nil
}

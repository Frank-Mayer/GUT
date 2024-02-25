package git

import (
	"errors"
	"os/exec"

	"github.com/go-git/go-git/v5"
)

var (
	Worktree     *git.Worktree
	Repo         *git.Repository
	gitInstalled bool
	gitBin       string
)

func Init() error {
	var err error
	gitBin, err = exec.LookPath("git")
	if err != nil || gitBin == "" {
		// git not found
		gitInstalled = false
	} else {
		// git found
		gitInstalled = true
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

	// check if there are any staged files
	status, err := Worktree.Status()
	if err != nil {
		return errors.Join(errors.New("failed to get status"), err)
	}
	if status.IsClean() {
		return errors.New("no local changes to commit")
	}

	return nil
}

package git

import (
	"errors"
	"os"
	"os/exec"

	"github.com/Frank-Mayer/gut/internal/config"
	"github.com/go-git/go-git/v5"
)

func Commit(message string) error {
	if gitInstalled {
		return cliCommit(message)
	} else {
		return goCommit(message)
	}
}

func cliCommit(message string) error {
	cmd := exec.Command(gitBin, "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return errors.Join(errors.New("failed to commit"), err)
	}
	return nil
}

func goCommit(message string) error {
	if _, err := Worktree.Commit(message, &git.CommitOptions{AllowEmptyCommits: config.Force}); err != nil {
		return errors.Join(errors.New("failed to commit"), err)
	}

	return nil
}

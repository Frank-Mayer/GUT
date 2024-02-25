package git

import (
	"errors"
	"os"
	"os/exec"
)

func Commit(message string) error {
	// check if there are any staged files
	status, err := Worktree.Status()
	if err != nil {
		return errors.Join(errors.New("failed to get status"), err)
	}
	if status.IsClean() {
		return errors.New("no local changes to commit")
	}

	if GitInstalled {
		return cliCommit(message)
	} else {
		return goCommit(message)
	}
}

func cliCommit(message string) error {
	cmd := exec.Command(GitBin, "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return errors.Join(errors.New("failed to commit"), err)
	}
	return nil
}

func goCommit(message string) error {
	if _, err := Worktree.Commit(message, nil); err != nil {
		return errors.Join(errors.New("failed to commit"), err)
	}

	return nil
}

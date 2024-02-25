package git

import (
	"errors"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
)

func Push() error {
	if GitInstalled {
		return cliPush()
	} else {
		return goPush()
	}
}

func cliPush() error {
	cmd := exec.Command(GitBin, "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return errors.Join(errors.New("failed to push"), err)
	}
	return nil
}

func goPush() error {
	if err := Repo.Push(&git.PushOptions{}); err != nil {
		return errors.Join(errors.New("failed to push"), err)
	}
	return nil
}

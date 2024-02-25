package git

import (
	"errors"
	"os"
	"os/exec"
)

func Stash(name string) error {
	if GitInstalled {
		return cliStash(name)
	} else {
		return goStash(name)
	}
}

func cliStash(name string) error {
	cmd := exec.Command(GitBin, "stash", "push", "-m", name, "--include-untracked")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func goStash(_ string) error {
	return errors.New("stash is not available in go-git")
}

func StashPop() error {
	if GitInstalled {
		return cliStashPop()
	} else {
		return goStashPop()
	}
}

func cliStashPop() error {
	cmd := exec.Command(GitBin, "stash", "pop", "--index")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		cmd := exec.Command(GitBin, "stash", "pop")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}
	return nil
}

func goStashPop() error {
	return errors.New("stash pop is not available in go-git")
}

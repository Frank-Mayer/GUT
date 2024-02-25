package git

import (
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
)

func Pull() error {
	if GitInstalled {
		return cliPull()
	} else {
		return goPull()
	}
}

func cliPull() error {
	cmd := exec.Command(GitBin, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func goPull() error {
	return Worktree.Pull(&git.PullOptions{
		Progress: os.Stdout,
	})
}

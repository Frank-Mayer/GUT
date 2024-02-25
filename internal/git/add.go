package git

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/go-git/go-git/v5"
)

// Add starts a interactive git add
func Add() error {
	// Find all changed files that can be added
	files, err := Worktree.Status()
	if err != nil {
		return err
	}
	// orde files alphabetically
	options := make([]huh.Option[string], 0, len(files))
	for file, status := range files {
		if status.Staging != git.Untracked && status.Staging != git.Unmodified { // is the file staged?
			options = append(options, huh.NewOption(file, file).Selected(true))
		} else if status.Worktree != git.Unmodified { // is the file modified?
			options = append(options, huh.NewOption(file, file))
		}
	}

	if len(options) == 0 {
		return errors.New("no local changes to commit")
	}

	var stagedFiles []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(options...).
				Value(&stagedFiles).
				Title("Select the files you want to add"),
		),
	)

	if err := form.Run(); err != nil {
		return errors.Join(errors.New("failed to get user input"), err)
	}

	if len(stagedFiles) == 0 {
		return errors.New("no files selected")
	}

	for _, file := range stagedFiles {
		if _, err := Worktree.Add(file); err != nil {
			return errors.Join(errors.New("failed to add file to the index"), err)
		}
	}

	return nil
}

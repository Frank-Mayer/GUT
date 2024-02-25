package cmd

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Frank-Mayer/gut/internal/git"
	"github.com/charmbracelet/huh"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file contents to the index",
	Long:  "This command updates the index using the current content found in the working tree, to prepare the content staged for the next commit. It typically adds the current content of existing paths as a whole, but with some options it can also be used to add content with only part of the changes made to the working tree files applied, or remove paths that do not exist in the working tree anymore.",
	RunE: func(_ *cobra.Command, args []string) error {
		if err := git.Init(); err != nil {
			return err
		}

		if len(args) != 0 {
			for _, arg := range args {
				if err := git.Worktree.AddGlob(arg); err != nil {
					return errors.Join(fmt.Errorf("Error adding %s to the index", arg), err)
				}
			}
			return nil
		}

		// interactive mode
		s, err := git.Worktree.Status()
		if err != nil {
			return errors.Join(fmt.Errorf("Error getting worktree status"), err)
		}

		// find all files that have local changes
		var files []string
		filesToSelect := new([]string)
		for file, status := range s {
			if status.Worktree != gogit.Unmodified {
				files = append(files, file)
			} else if status.Staging == gogit.Added || status.Staging == gogit.Modified || status.Staging == gogit.Deleted || status.Staging == gogit.Renamed {
				*filesToSelect = append(*filesToSelect, file)
				files = append(files, file)
			}
		}
		// sort files
		sort.Strings(files)

		form := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select files to add").
				Options(huh.NewOptions[string](files...)...).
				Value(filesToSelect),
		))
		if err := form.Run(); err != nil {
			return err
		}
		for _, file := range *filesToSelect {
			if _, err := git.Worktree.Add(file); err != nil {
				return errors.Join(fmt.Errorf("Error adding %s to the index", file), err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

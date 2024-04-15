package cmd

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/tsukinoko-kun/gut/internal/git"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls the latest changes from the remote repository.",
	Long: "Pulls the latest changes from the remote repository. " +
		"If you have local changes, they will be merged with the remote changes.",
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := git.Init(); err != nil {
			return err
		}
		s, err := git.Worktree.Status()
		if err != nil {
			return errors.Join(fmt.Errorf("failed to get worktree status"), err)
		}

		stashed := false
		var stashName string
		if !s.IsClean() {
			fmt.Println("You have local changes.")
			// stash changes
			stashName = fmt.Sprintf("gut-temp-%x", rand.Int63())
			if err := git.Stash(stashName); err != nil {
				return errors.Join(fmt.Errorf("failed to stash local changes"), err)
			}
			stashed = true
		}

		// pull changes
		if err := git.Pull(); err != nil {
			if stashed {
				// pop stash
				if err := git.StashPop(); err != nil {
					fmt.Printf("Failed to pop stash. Please do it manually: %s\n", stashName)
					fmt.Println(err)
				}
			}

			return errors.Join(fmt.Errorf("failed to pull changes"), err)
		}

		if stashed {
			// pop stash
			if err := git.StashPop(); err != nil {
				return errors.Join(fmt.Errorf("failed to pop stash"), err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

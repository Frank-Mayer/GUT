package cmd

import (
	"errors"

	"github.com/Frank-Mayer/gut/internal/git"
	"github.com/Frank-Mayer/gut/internal/input"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var (
	commitInput   = input.CommitInput{}
	commitBreak   bool
	commitNoBreak bool

	// commitCmd represents the commit command
	commitCmd = &cobra.Command{
		Use:   "commit",
		Short: "Create a conventional commit",
		Long: `This command will guide you through the process of creating a conventional commit.
You will be asked for the type, scope, description, body and footer of the commit.
You need to stage the changes before running this command.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := git.Init(); err != nil {
				return err
			}

			if s, err := git.Worktree.Status(); err != nil {
				return err
			} else {
				// check for staged files
				hasStagedFiles := false
				for _, status := range s {
					if status.Staging == gogit.Modified || status.Staging == gogit.Added || status.Staging == gogit.Deleted || status.Staging == gogit.Renamed || status.Staging == gogit.Copied || status.Staging == gogit.UpdatedButUnmerged {
						hasStagedFiles = true
						break
					}
				}
				if !hasStagedFiles {
					return errors.New("You need to stage changes before you can commit.")
				}
			}
			if commitBreak && commitNoBreak {
				return errors.New("You can't use both --break and --no-break")
			}

			if commitBreak {
				commitInput.BreakingChanges = input.True
			} else if commitNoBreak {
				commitInput.BreakingChanges = input.False
			} else {
				commitInput.BreakingChanges = input.Undefined
			}

			if err := commitInput.Ask(); err != nil {
				return err
			}

			cMessage, err := commitInput.String()
			if err != nil {
				return err
			}

			c, err := input.Choose(cMessage, "commit", "commit and push", "abort")
			if err != nil {
				return err
			}

			switch c {
			case 0:
				if err := git.Commit(cMessage); err != nil {
					return err
				}
			case 1:
				if err := git.Commit(cMessage); err != nil {
					return err
				}
				if err := git.Push(); err != nil {
					return err
				}
			case 2:
				return nil
			default:
				return errors.New("Invalid choice")
			}
			return nil
		},
	}
)

func init() {
	commitCmd.Flags().StringVarP(&commitInput.Type, "type", "t", "", "The type of the commit")
	commitCmd.Flags().BoolVarP(&commitBreak, "break", "b", false, "Add a line break after the type")
	commitCmd.Flags().BoolVarP(&commitNoBreak, "no-break", "B", false, "Do not add a line break after the type")
	commitCmd.Flags().StringVarP(&commitInput.Scope, "scope", "s", "", "The scope of the commit")
	commitCmd.Flags().StringVarP(&commitInput.Description, "description", "d", "", "The description of the commit")
	commitCmd.Flags().StringVar(&commitInput.Body, "body", "", "The body of the commit")
	commitCmd.Flags().StringVar(&commitInput.Footer, "footer", "", "The footer of the commit")

	rootCmd.AddCommand(commitCmd)
}

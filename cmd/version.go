package cmd

import (
	"github.com/Frank-Mayer/gut/internal/git"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the the current version from tags",
	Long:  "Looking for the latest tag containing a semantic version, and print it to the standard output",
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := git.Init(); err != nil {
			return err
		}
		v, _, err := git.SemverStatus()
		if err != nil {
			return err
		}
		if v != nil {
			println(v.String())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

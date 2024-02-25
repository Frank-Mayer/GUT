package cmd

import (
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty Git repository or reinitialize an existing one",
	Long:  `This command creates an empty Git repository - basically a .git directory with subdirectories for objects, refs/heads, refs/tags, and template files. An initial HEAD file that references the HEAD of the master branch is also created.`,
	RunE: func(c *cobra.Command, _ []string) error {
		if gitBin, err := exec.LookPath("git"); err == nil && gitBin != "" {
			cmd := exec.Command(gitBin, os.Args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				return err
			}
		} else {
			isBare, _ := c.Flags().GetBool("bare")
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			if _, err := git.PlainInit(wd, isBare); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	initCmd.Flags().BoolP("quiet", "q", false, "Only print error and warning messages; all other output will be suppressed.")
	initCmd.Flags().Bool("bare", false, "Create a bare repository. If `GIT_DIR` environment is not set, it is set to the current working directory.")
	initCmd.Flags().String("object-format", "sha1", "Specify the given object format (hash algorithm) for the repository. The valid values are sha1 and (if enabled) sha256. sha1 is the default. Note: At present, there is no interoperability between SHA-256 repositories and SHA-1 repositories.")
	initCmd.Flags().String("ref-format", "files", "Specify the given ref storage format for the repository.")
	initCmd.Flags().String("template", "", "Specify the directory from which templates will be used. (See the 'TEMPLATE DIRECTORY' section of git-init(1).)")
	initCmd.Flags().String("separate-git-dir", ".git", "Instead of initializing the repository as a directory to either `$GIT_DIR` or `./.git/`, create a text file there containing the path to the actual repository. This file acts as a filesystem-agnostic Git symbolic link to the repository. If this is a reinitialization, the repository will be moved to the specified path.")
	initCmd.Flags().StringP("initial-branch", "b", "main", "Use the specified name for the initial branch in the newly created repository. If not specified, fall back to the default name (main).")
	initCmd.Flags().String("shared", "umask", "Specify that the Git repository is to be shared amongst several users. This allows users belonging to the same group to push into that repository. When specified, the config variable \"core.sharedRepository\" is set so that files and directories under `$GIT_DIR` are created with the requested permissions. When not specified, Git will use permissions reported by umask(2).")

	rootCmd.AddCommand(initCmd)
}

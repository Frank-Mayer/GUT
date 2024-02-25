package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Frank-Mayer/gut/internal/git"
	"github.com/Frank-Mayer/gut/internal/semver"
	"github.com/spf13/cobra"
)

type log struct {
	CurrentVersion string           `json:"current_version"`
	NextVersion    string           `json:"next_version"`
	Changes        []*semver.Change `json:"changes"`
}

var logFormat string = "json"

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Get the change log since the last release",
	Long:  "Get the change log since the last release.",
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := git.Init(); err != nil {
			return err
		}

		v, c, err := git.SemverStatus()
		if err != nil {
			return err
		}

		var l log

		l.CurrentVersion = v.String()
		l.Changes = c
		bump := semver.NoBump
		for _, change := range c {
			if change.Bump > bump {
				bump = change.Bump
			}
		}
		l.NextVersion = v.Bump(bump).String()

		switch logFormat {
		case "json":
			return printJSON(l)
		case "markdown", "md":
			return printMarkdown(l)
		default:
			return fmt.Errorf("invalid format: %q", logFormat)
		}
	},
}

func init() {
	logCmd.Flags().StringVarP(&logFormat, "format", "f", logFormat, "Output format (json, markdown)")
	rootCmd.AddCommand(logCmd)
}

func printJSON(l interface{}) error {
	str, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}
	fmt.Println(string(str))
	return nil
}

func printMarkdown(l log) error {
	fmt.Printf("## %s\n\n", l.NextVersion)
	fmt.Printf("*From %s*\n\n", l.CurrentVersion)

	hasBreakingChanges := false
	for _, change := range l.Changes {
		if change.Bump == semver.MajorBump {
			hasBreakingChanges = true
			break
		}
	}
	if hasBreakingChanges {
		fmt.Print("### Breaking Changes\n\n")
		for _, change := range l.Changes {
			if change.Bump == semver.MajorBump {
				printMarkdownChange(change)
			}
		}
		fmt.Print("### All Changes\n\n")
	} else {
		fmt.Print("### Changes\n\n")
	}

	for _, change := range l.Changes {
		printMarkdownChange(change)
	}
	return nil
}

func printMarkdownChange(change *semver.Change) {
	fmt.Print("- ")
	if change.Scope != "" {
		fmt.Printf("**%s**: ", change.Scope)
	}
	fmt.Print(change.Description)
	if change.Body == "" {
		fmt.Print("\n")
	} else {
		fmt.Print("  \n")
		lines := strings.Split(change.Body, "\n")
		for i := 0; i < len(lines); i++ {
			l := strings.TrimSpace(lines[i])
			lastLine := i == len(lines)-1
			if l != "" {
				fmt.Printf("  %s", l)
			}
			if lastLine {
				fmt.Print("\n")
			} else {
				fmt.Print("  \n")
			}
		}
	}
}

package input

import (
	"errors"
	"os"
	"regexp"

	"github.com/charmbracelet/huh"
)

func (i *CommitInput) Ask() error {
	accessibleMode := os.Getenv("ACCESSIBLE") != ""

	fields := make([]huh.Field, 0)

	if i.Type == "" {
		fields = append(fields,
			huh.NewSelect[string]().
				Title("Type").
				Options(
					huh.NewOption("Feature (introduces a new feature)", "feat"),
					huh.NewOption("Fix (patches a bug)", "fix"),
					huh.NewOption("Build (changes that affect the build system or external dependencies)", "build"),
					huh.NewOption("Chore (other changes that don't modify src or test files)", "chore"),
					huh.NewOption("CI (changes to our CI configuration files and scripts)", "ci"),
					huh.NewOption("Docs (documentation only changes)", "docs"),
					huh.NewOption("Style (changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc))", "style"),
					huh.NewOption("Refactor (a code change that neither fixes a bug nor adds a feature)", "refactor"),
					huh.NewOption("Performance (a code change that improves performance without adding a feature or fixing a bug)", "perf"),
					huh.NewOption("Test (adding missing tests or correcting existing tests)", "test"),
				).
				Value(&i.Type),
		)
	}

	if i.Scope == "" {
		fields = append(fields,
			huh.NewInput().
				Title("Scope").
				Value(&i.Scope),
		)
	}

	if i.BreakingChanges == Undefined {
		fields = append(fields,
			huh.NewSelect[Bool]().
				Title("Breaking Changes").
				Options(
					huh.NewOption("Yes", True),
					huh.NewOption("No", False),
				).
				Value(&i.BreakingChanges),
		)
	}

	if i.Description == "" {
		fields = append(fields,
			huh.NewInput().
				Title("Description").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("description cannot be empty")
					}
					return nil
				}).
				Value(&i.Description),
		)
	}

	if len(fields) != 0 {
		if err := huh.NewForm(huh.NewGroup(fields...)).
			WithAccessible(accessibleMode).
			Run(); err != nil {
			return errors.Join(errors.New("faild to get user input"), err)
		}
		fields = make([]huh.Field, 0)
	}

	if i.Body == "" {
		fields = append(fields,
			huh.NewText().
				Title("Body").
				Value(&i.Body),
		)
	}

	if i.Footer == "" {
		// /^breaking[\s_-]*changes?\s*[:=]?\s*\w+/im
		re := regexp.MustCompile(`(?i)breaking[\s_-]*changes?\s*[:=]?\s*\w+`)
		if i.BreakingChanges == True {
			i.Footer = "BREAKING CHANGE: "
		}
		fields = append(fields,
			huh.NewText().
				Title("Footer").
				Value(&i.Footer).
				Validate(func(s string) error {
					containsNote := re.MatchString(s)
					if i.BreakingChanges == True {
						if !containsNote {
							return errors.New("footer must contain a breaking change note `BREAKING CHANGE: ...`")
						}
					} else if containsNote {
						return errors.New("footer must not contain a breaking change note `BREAKING CHANGE: ...`")
					}
					return nil
				}),
		)
	}

	if len(fields) > 0 {
		s, _ := i.String()
		fields = append(fields, huh.NewNote())
		if err := huh.NewForm(huh.NewGroup(fields...).Description(s)).
			WithAccessible(accessibleMode).
			Run(); err != nil {
			return errors.Join(errors.New("faild to get user input"), err)
		}

	}

	return nil
}

func Confirm(s string) bool {
	ok := false
	form := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Value(&ok).Title(s),
	))
	if err := form.Run(); err != nil {
		return false
	}
	return ok
}

func Choose(s string, options ...string) (int, error) {
	var choice int
	huhOptions := make([]huh.Option[int], 0, len(options))
	for i, option := range options {
		huhOptions = append(huhOptions, huh.NewOption(option, i))
	}
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[int]().
			Title(s).
			Options(huhOptions...).
			Value(&choice),
	))
	if err := form.Run(); err != nil {
		return -1, errors.Join(errors.New("faild to get user input"), err)
	}
	return choice, nil
}

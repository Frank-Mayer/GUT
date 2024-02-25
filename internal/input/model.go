package input

import (
	"fmt"
	"strings"
)

type Bool uint8

const (
	Undefined Bool = iota
	False
	True
)

// validation
var (
	validTypes = []string{"feat", "fix", "build", "chore", "ci", "docs", "style", "refactor", "perf", "test"}
)

type CommitInput struct {
	Type            string
	BreakingChanges Bool
	Scope           string
	Description     string
	Body            string
	Footer          string
}

func (i *CommitInput) ValidateType() error {
	if i.Type == "" {
		return nil
	}

	for _, t := range validTypes {
		if t == i.Type {
			return nil
		}
	}

	return fmt.Errorf("invalid type: %s", i.Type)

}

func (i CommitInput) String() (string, error) {
	var s strings.Builder

	i.Type = strings.TrimSpace(i.Type)
	if i.Type == "" {
		return "", fmt.Errorf("type cannot be empty")
	}
	if err := i.ValidateType(); err != nil {
		return "", err
	}
	s.WriteString(i.Type)

	i.Scope = strings.TrimSpace(i.Scope)
	if i.Scope != "" {
		s.WriteRune('(')
		s.WriteString(i.Scope)
		s.WriteRune(')')
	}

	if i.BreakingChanges == True {
		s.WriteRune('!')
	}

	s.WriteString(": ")
	i.Description = strings.TrimSpace(i.Description)
	if i.Description == "" {
		return "", fmt.Errorf("description cannot be empty")
	}
	s.WriteString(i.Description)

	i.Body = strings.TrimSpace(i.Body)
	if i.Body != "" {
		s.WriteString("\n\n")
		s.WriteString(i.Body)
	}

	i.Footer = strings.TrimSpace(i.Footer)
	if i.Footer != "" {
		s.WriteString("\n\n")
		s.WriteString(i.Footer)
	}

	return s.String(), nil
}

package semver

import (
	"regexp"
	"strings"
)

type (
	BumpType uint8

	Change struct {
		Type        string
		Scope       string
		Description string
		Body        string
		Bump        BumpType
	}
)

const (
	NoBump BumpType = iota
	PatchBump
	MinorBump
	MajorBump
)

var commitMessageRegex *regexp.Regexp

func CommitMessageRegex() *regexp.Regexp {
	if commitMessageRegex == nil {
		commitMessageRegex = regexp.MustCompile(`(\w+)\s*(\(\s*([^)]+)\s*\))?\s*(!?)\s*:\s*(.+)\s*`)
	}
	return commitMessageRegex
}

func ChangeFromCommitMessage(message string) *Change {
	lines := strings.Split(message, "\n")
	firstLine := lines[0]

	matches := CommitMessageRegex().FindStringSubmatch(firstLine)
	if len(matches) != 6 {
		return &Change{
			Description: firstLine,
			Body:        strings.TrimSpace(strings.Join(lines[1:], "\n")),
		}
	}

	cType := matches[1]
	cScope := matches[3]
	cBreaking := matches[4]
	cDescription := matches[5]

	c := new(Change)
	c.Type = cType

	lct := strings.ToLower(cType)
	switch lct {
	case "feat", "feature":
		c.Type = "feat"
		c.Bump = MinorBump
	case "fix", "bugfix":
		c.Type = "fix"
		c.Bump = PatchBump
	case "breaking", "major":
		c.Type = "breaking"
		c.Bump = MajorBump
	default:
		c.Type = lct
		c.Bump = NoBump
	}

	c.Scope = cScope
	c.Description = strings.TrimSpace(cDescription)
	c.Body = strings.TrimSpace(strings.Join(lines[1:], "\n"))

	if cBreaking == "!" {
		c.Bump = MajorBump
	}

	return c
}

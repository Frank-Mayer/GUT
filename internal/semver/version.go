package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func NewVersion(major, minor, patch uint64) *Version {
	v := new(Version)
	v.Major = major
	v.Minor = minor
	v.Patch = patch
	return v
}

var tagParseRegex *regexp.Regexp

func TagParseRegex() *regexp.Regexp {
	if tagParseRegex == nil {
		tagParseRegex = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)\s*$`)
	}
	return tagParseRegex
}

func ParseVersion(version string) (*Version, error) {
	v := new(Version)
	matches := TagParseRegex().FindStringSubmatch(version)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid version string")
	}
	var err error
	if v.Major, err = strconv.ParseUint(matches[1], 10, 64); err != nil {
		return nil, errors.Join(fmt.Errorf("invalid major version: %s", matches[1]), err)
	}
	if v.Minor, err = strconv.ParseUint(matches[2], 10, 64); err != nil {
		return nil, errors.Join(fmt.Errorf("invalid minor version: %s", matches[2]), err)
	}
	if v.Patch, err = strconv.ParseUint(matches[3], 10, 64); err != nil {
		return nil, errors.Join(fmt.Errorf("invalid patch version: %s", matches[3]), err)
	}
	return v, nil
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		} else {
			return -1
		}
	}
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		} else {
			return -1
		}
	}
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

func (v *Version) Equals(other *Version) bool {
	return v.Compare(other) == 0
}

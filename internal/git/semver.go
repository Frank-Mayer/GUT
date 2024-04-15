package git

import (
	"errors"
	"fmt"

	"github.com/tsukinoko-kun/gut/internal/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

func SemverStatus() (*semver.Version, []*semver.Change, error) {
	tags := make(map[string]string)
	if t, err := Repo.Tags(); err != nil {
		return nil, nil, errors.Join(errors.New("could not get tags"), err)
	} else {
		err := t.ForEach(func(t *plumbing.Reference) error {
			if t == nil {
				return nil
			}
			tag, err := Repo.TagObject(t.Hash())
			if err != nil {
				return errors.Join(fmt.Errorf("could not get tag object for tag %q", t.Name().String()), err)
			}
			commit, err := tag.Commit()
			if err != nil {
				return errors.Join(fmt.Errorf("could not get commit object for tag %q", t.Name().String()), err)
			}
			hash := commit.Hash.String()
			name := tag.Name
			tags[hash] = name
			return nil
		})
		if err != nil {
			return nil, nil, errors.Join(errors.New("could not get tags"), err)
		}
	}

	history, err := Repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, nil, errors.Join(errors.New("could not get worktree status"), err)
	}

	changes := make([]*semver.Change, 0)
	var latestTag *semver.Version

	err = history.ForEach(func(c *object.Commit) error {
		if c == nil {
			return nil
		}

		// check if commit is a tag
		hash := c.Hash.String()
		if tag, ok := tags[hash]; ok {
			latestTag, err = semver.ParseVersion(tag)
			if err != nil {
				return nil
			}
			return storer.ErrStop
		}

		// store changes
		change := semver.ChangeFromCommitMessage(c.Message)
		changes = append(changes, change)
		return nil
	})
	if err != nil {
		return nil, nil, errors.Join(errors.New("could not get worktree status"), err)
	}
	if latestTag == nil {
		return nil, nil, errors.New("no version tag found")
	}

	return latestTag, changes, nil
}

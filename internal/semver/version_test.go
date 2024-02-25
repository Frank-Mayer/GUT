package semver_test

import (
	"testing"

	"github.com/Frank-Mayer/gut/internal/semver"
)

func TestVersion_Compare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a    semver.Version
		b    semver.Version
		want int
	}{
		{semver.Version{1, 0, 0}, semver.Version{1, 0, 0}, 0},
		{semver.Version{1, 0, 0}, semver.Version{1, 0, 1}, -1},
		{semver.Version{1, 0, 0}, semver.Version{1, 1, 0}, -1},
		{semver.Version{1, 0, 0}, semver.Version{2, 0, 0}, -1},
		{semver.Version{1, 0, 0}, semver.Version{0, 0, 0}, 1},
		{semver.Version{1, 0, 0}, semver.Version{0, 0, 14}, 1},
		{semver.Version{1, 0, 32}, semver.Version{0, 26, 5}, 1},
		{semver.Version{1, 0, 0}, semver.Version{2, 0, 1}, -1},
		{semver.Version{1, 0, 0}, semver.Version{1, 1, 1}, -1},
		{semver.Version{1, 0, 0}, semver.Version{2, 1, 0}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.a.String()+"-"+tt.b.String(), func(t *testing.T) {
			t.Parallel()
			got := tt.a.Compare(&tt.b)
			if got != tt.want {
				t.Errorf("Version.Compare(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		txt  string
		want semver.Version
	}{
		{"1.0.0", semver.Version{1, 0, 0}},
		{"1.0.1", semver.Version{1, 0, 1}},
		{"1.1.0", semver.Version{1, 1, 0}},
		{"2.0.0", semver.Version{2, 0, 0}},
		{"0.0.0", semver.Version{0, 0, 0}},
		{"0.0.14", semver.Version{0, 0, 14}},
		{"0.26.5", semver.Version{0, 26, 5}},
		{"2.0.1", semver.Version{2, 0, 1}},
		{"1.1.1", semver.Version{1, 1, 1}},
		{"2.1.0", semver.Version{2, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.txt, func(t *testing.T) {
			t.Parallel()
			got, err := semver.ParseVersion(tt.txt)
			if err != nil {
				t.Errorf("ParseVersion(%v) returned error: %v", tt.txt, err)
			}
			if *got != tt.want {
				t.Errorf("ParseVersion(%v) = %v; want %v", tt.txt, got, tt.want)
			}
		})
	}
}

func TestParseVersionError(t *testing.T) {
	ftests := []struct {
		txt string
	}{
		{""},
		{"1.0"},
		{" 1.0.0-rc1"},
		{"1.0.0-rc1 "},
		{"1.0.0-rc1"},
		{"foo"},
	}
	for _, tt := range ftests {
		t.Run(tt.txt, func(t *testing.T) {
			t.Parallel()
			_, err := semver.ParseVersion(tt.txt)
			if err == nil {
				t.Errorf("ParseVersion(%v) did not return error", tt.txt)
			}
		})
	}
}

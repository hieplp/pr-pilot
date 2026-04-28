package git_test

import (
	"strings"
	"testing"

	"github.com/hieplp/pr-pilot/internal/git"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name      string
		diff      string
		maxBytes  int
		wantCut   bool
		wantBegin string
	}{
		{"below limit", "small diff", 100, false, "small diff"},
		{"exactly at limit", "exactly", 7, false, "exactly"},
		{"above limit", "this is a long diff", 5, true, "this "},
		{"zero disables truncation", "any content", 0, false, "any content"},
		{"negative disables truncation", "any content", -1, false, "any content"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := git.Truncate(tt.diff, tt.maxBytes)

			hasTruncNotice := strings.Contains(got, "[diff truncated")
			if hasTruncNotice != tt.wantCut {
				t.Errorf("truncation notice present=%v, want %v", hasTruncNotice, tt.wantCut)
			}
			if !strings.HasPrefix(got, tt.wantBegin) {
				t.Errorf("result starts with %q, want prefix %q", got, tt.wantBegin)
			}
		})
	}
}

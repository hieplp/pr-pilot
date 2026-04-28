package prompt_test

import (
	"strings"
	"testing"

	"github.com/hieplp/pr-pilot/internal/prompt"
)

func TestCommitPrompt(t *testing.T) {
	system, user := prompt.CommitPrompt("diff --git a/foo.go")

	if !strings.Contains(system, "Conventional Commits") {
		t.Error("system prompt should mention Conventional Commits")
	}
	if !strings.Contains(user, "diff --git a/foo.go") {
		t.Error("user message should contain the diff")
	}
	if !strings.Contains(user, "```diff") {
		t.Error("user message should wrap diff in a code fence")
	}
}

func TestPRPrompt(t *testing.T) {
	system, user := prompt.PRPrompt("feature/auth", "main", "diff content", "abc1234 add login")

	if !strings.Contains(strings.ToLower(system), "pull request") {
		t.Error("system prompt should mention pull request")
	}
	for _, want := range []string{"feature/auth", "main", "diff content", "abc1234 add login"} {
		if !strings.Contains(user, want) {
			t.Errorf("user message should contain %q", want)
		}
	}
}

func TestPRTitle(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"first content line after summary heading", "## Summary\nAdd user auth\n\n## Motivation\nNeeded for login.", "Add user auth"},
		{"plain first line", "Add user auth\n\nMore details", "Add user auth"},
		{"empty body", "", "PR description"},
		{"only heading with no content", "## Summary\n", "PR description"},
		{"heading markers stripped", "### Add feature", "Add feature"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prompt.PRTitle(tt.body)
			if got != tt.want {
				t.Errorf("PRTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

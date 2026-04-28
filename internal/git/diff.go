package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// StagedDiff returns the output of `git diff --cached`.
func StagedDiff() (string, error) {
	out, err := run("git", "diff", "--cached")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(out) == "" {
		return "", errors.New("no staged changes — run `git add` first")
	}
	return out, nil
}

// BranchDiff returns the diff between base and HEAD (git diff <base>...HEAD).
func BranchDiff(base string) (string, error) {
	out, err := run("git", "diff", base+"...HEAD")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(out) == "" {
		return "", errors.New("no changes between current branch and " + base)
	}
	return out, nil
}

// CommitLog returns the one-line commit log between base and HEAD.
func CommitLog(base string) (string, error) {
	return run("git", "log", "--oneline", base+"..HEAD")
}

// CurrentBranch returns the name of the current branch.
func CurrentBranch() (string, error) {
	return run("git", "rev-parse", "--abbrev-ref", "HEAD")
}

// Commit runs `git commit -m <message>`.
func Commit(message string) error {
	_, err := run("git", "commit", "-m", message)
	return err
}

// PRTemplate returns the contents of the repo's PR template file if one exists.
// Checks common locations used by GitHub: .github/, docs/, and repo root.
func PRTemplate() string {
	candidates := []string{
		".github/pull_request_template.md",
		".github/PULL_REQUEST_TEMPLATE.md",
		"docs/pull_request_template.md",
		"pull_request_template.md",
	}
	for _, p := range candidates {
		if b, err := os.ReadFile(p); err == nil {
			return string(b)
		}
	}
	return ""
}

// Truncate caps diff to maxBytes and appends a notice when it was cut.
// A maxBytes value of 0 or less disables truncation.
func Truncate(diff string, maxBytes int) string {
	if maxBytes <= 0 || len(diff) <= maxBytes {
		return diff
	}
	return diff[:maxBytes] + fmt.Sprintf("\n[diff truncated — showing first %d bytes]", maxBytes)
}

// Push pushes the given branch to origin, setting the upstream if not already set.
func Push(branch string) error {
	_, err := run("git", "push", "-u", "origin", branch)
	return err
}

func run(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New(msg)
	}
	return stdout.String(), nil
}

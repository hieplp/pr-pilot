package git

import (
	"bytes"
	"errors"
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

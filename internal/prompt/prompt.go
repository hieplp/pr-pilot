package prompt

import "fmt"

const commitSystem = `You are an expert software engineer writing git commit messages.
Follow the Conventional Commits specification strictly:
  <type>(<optional scope>): <short summary>

  [optional body]

  [optional footer(s)]

Allowed types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert.
Rules:
- Summary line ≤ 72 characters, imperative mood, no period at end.
- Body explains *why*, not *what* (the diff already shows what).
- Add BREAKING CHANGE footer when applicable.
- Output ONLY the commit message — no explanation, no markdown fences.`

const prSystem = `You are an expert software engineer writing pull request descriptions.
Produce a structured PR body in Markdown with these sections:
  ## Summary
  ## Motivation
  ## Changes
  ## Test Plan

Rules:
- Be concise and specific — reference actual file or function names from the diff.
- Use bullet points inside each section.
- Output ONLY the PR description — no explanation, no markdown fences around the whole thing.`

// CommitPrompt returns the system instruction and user message for commit generation.
func CommitPrompt(diff string) (system, user string) {
	return commitSystem, fmt.Sprintf("Generate a commit message for the following staged diff:\n\n```diff\n%s\n```", diff)
}

// PRPrompt returns the system instruction and user message for PR description generation.
func PRPrompt(branch, base, diff, log string) (system, user string) {
	return prSystem, fmt.Sprintf(
		"Branch: %s → %s\n\nCommit log:\n%s\n\nDiff:\n```diff\n%s\n```",
		branch, base, log, diff,
	)
}

// PRTitle extracts a short PR title from the first non-empty line of a generated PR body.
func PRTitle(body string) string {
	for _, line := range splitLines(body) {
		if line != "" && line != "## Summary" {
			// strip leading markdown heading markers
			for len(line) > 0 && (line[0] == '#' || line[0] == ' ') {
				line = line[1:]
			}
			if line != "" {
				return line
			}
		}
	}
	return "PR description"
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

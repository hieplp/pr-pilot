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

// CommitPrompt builds the full prompt for commit message generation.
func CommitPrompt(diff string) string {
	return fmt.Sprintf("%s\n\nGenerate a commit message for the following staged diff:\n\n```diff\n%s\n```", commitSystem, diff)
}

// PRPrompt builds the full prompt for PR description generation.
func PRPrompt(branch, base, diff, log string) string {
	return fmt.Sprintf(
		"%s\n\nBranch: %s → %s\n\nCommit log:\n%s\n\nDiff:\n```diff\n%s\n```",
		prSystem, branch, base, log, diff,
	)
}

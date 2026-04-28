# pr-pilot Implementation Progress

## Phase 1 — `commit` end-to-end (in progress)

| # | File | Status | Notes |
|---|------|--------|-------|
| 1 | `internal/config/config.go` | [x] | Viper setup: flags → env vars → TOML file |
| 2 | `internal/git/diff.go` | [x] | `StagedDiff()` and `BranchDiff(base)` via `os/exec` |
| 3 | `internal/prompt/prompt.go` | [x] | Conventional Commits prompt template |
| 4 | `internal/provider/claude.go` | [x] | Anthropic SDK implementation |
| 5 | `internal/provider/factory.go` | [x] | Maps `--provider` flag → concrete `Provider` |
| 6 | `cmd/commit.go` (wire) | [x] | flags → config → diff → prompt → provider → print |

## Phase 2 — Provider coverage + `pr` command

| # | File | Status | Notes |
|---|------|--------|-------|
| 7 | `internal/provider/openai.go` | [x] | OpenAI + Ollama (same client, different base URL) |
| 8 | `internal/prompt/prompt.go` | [x] | PR description prompt template (added to existing file) |
| 9 | `cmd/pr.go` (wire) | [x] | Reuse git/prompt/provider, PR-specific prompt |

## Phase 3 — Polish

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 10 | `.pr-pilot.toml` project config | [x] | Per-repo provider/model/base-branch overrides; layered on global config |
| 11 | Interactive TUI review | [x] | bubbletea: accept / edit ($EDITOR) / regenerate / quit; `--yes` skips |
| 12 | `internal/store/store.go` | [ ] | SQLite usage log (timestamp, tokens, cost) |
| 13 | `cmd/hook.go` | [x] | `pr-pilot hook install/uninstall` (prepare-commit-msg) |
| 14 | Shell completions | [x] | Built into Cobra — `pr-pilot completion bash/zsh/fish/powershell` |
| 15 | `--dry-run` token/cost estimate | [ ] | Show token count + estimated cost before sending |

## Roadmap (not yet scoped)

- Direct GitHub/GitLab PR creation
- Changelog generation (`pr-pilot changelog`)
- Branch name suggestions
- Hallucination guard (verify filenames in diff)
- Retry + fallback chain across providers
- Diff filtering (exclude lock files, generated files)

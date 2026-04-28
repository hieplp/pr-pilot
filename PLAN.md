# pr-pilot — Implementation Plan

## Current State

The core pipeline is fully wired:
`staged diff / branch diff → prompt builder → LLM provider → TUI review (accept / edit / regenerate / quit)`

Commands shipped: `commit`, `pr`, `hook install/uninstall`, `config show/model`  
Providers shipped: Claude (Anthropic SDK), OpenAI, Ollama (OpenAI-compatible)

---

## Roadmap

### Phase 2 — Close the Loop

#### 2.1 GitHub PR creation (`pr --create`)

**Goal:** `pr-pilot pr --create` generates the description *and* opens the PR on GitHub.

**Work:**
- Add `--create` flag to `cmd/pr.go`
- After `ActionAccept`, extract a short title from the first line of the generated description
- Shell out to `gh pr create --title "<title>" --body "<body>"` (same pattern as `git.Commit`)
- Add `git.Push(branch string)` in `internal/git/diff.go` for `git push -u origin <branch>`
- Add optional `--push` flag that runs push before create
- Error clearly if `gh` is not installed or the user is not authenticated

**Files:** `cmd/pr.go`, `internal/git/diff.go`

---

#### 2.2 Fix system/user message separation in Claude provider

**Goal:** Use the Anthropic API's dedicated `System` field instead of prepending the system prompt into the user message.

**Work:**
- Split `CommitPrompt` / `PRPrompt` into `(system string, user string)` return values
- Update `Provider` interface: `Complete(ctx, system, user string) (string, error)`
- Update `ClaudeProvider.Complete` to pass `System: system` in `MessageNewParams`
- Update `OpenAIProvider.Complete` to prepend a `system` role message
- Update both call sites in `cmd/commit.go` and `cmd/pr.go`

**Files:** `internal/prompt/prompt.go`, `internal/provider/provider.go`, `internal/provider/claude.go`, `internal/provider/openai.go`, `cmd/commit.go`, `cmd/pr.go`

---

### Phase 3 — UX Polish

#### 3.1 Streaming output + loading spinner

**Goal:** Show a spinner (or stream tokens) while waiting for the LLM, eliminating silent wait.

**Work:**
- Add a `bubbletea` spinner model that runs while `p.Complete()` is in flight
- Run the LLM call in a goroutine; send result back via `tea.Cmd`
- Alternatively, implement streaming in `ClaudeProvider` using `client.Messages.NewStreaming` and print tokens as they arrive before handing off to TUI review

**Files:** `internal/tui/review.go` (new spinner model), `internal/provider/claude.go`, `internal/provider/openai.go`

---

#### 3.2 Copy to clipboard on accept

**Goal:** Add `[c]` keybinding in the review TUI to copy generated content to the clipboard.

**Work:**
- In `tui/review.go`, add `ActionCopy` constant
- Handle `"c"` key in `Update()`: call `clipboard.WriteAll(content)` from `atotto/clipboard` (already in `go.mod`)
- Show a brief confirmation message in `View()`

**Files:** `internal/tui/review.go`

---

### Phase 4 — Correctness & Reliability

#### 4.1 Diff truncation

**Goal:** Prevent silent API errors when diffs exceed model context limits.

**Work:**
- Add `const maxDiffBytes = 80_000` in `internal/git/diff.go` (roughly ~20k tokens)
- If diff exceeds limit, truncate and append a warning line: `\n[diff truncated — showing first 80 000 bytes]`
- Make the limit overridable via config key `max_diff_bytes`

**Files:** `internal/git/diff.go`, `internal/config/config.go`, `internal/config/load.go`

---

#### 4.2 Configurable Ollama base URL

**Goal:** Allow users running Ollama on a non-default host/port to configure it.

**Work:**
- Add `OllamaBaseURL string` to `Config` struct
- Add `viper.SetDefault("ollama_base_url", "http://localhost:11434/v1")` in `load.go`
- Pass `cfg.OllamaBaseURL` through `provider.New()` to `NewOllama(model, baseURL string)`

**Files:** `internal/config/config.go`, `internal/config/load.go`, `internal/provider/openai.go`, `internal/provider/factory.go`, `cmd/commit.go`, `cmd/pr.go`

---

### Phase 5 — Developer Ergonomics

#### 5.1 Unit tests

**Goal:** Establish a test baseline. Target ≥ 80% coverage on pure-logic packages.

**Scope:**
- `internal/prompt` — table-driven tests for `CommitPrompt` / `PRPrompt` output shape
- `internal/config` — `Load()` with env-var overrides, `Override()` precedence
- `internal/provider` — stub `Provider` for use in command tests
- `cmd/commit`, `cmd/pr` — integration tests using the stub provider and a temp git repo

**Files:** `internal/prompt/prompt_test.go`, `internal/config/load_test.go`, `internal/provider/stub.go`, `cmd/commit_test.go`, `cmd/pr_test.go`

---

#### 5.2 `config set <key> <value>`

**Goal:** Allow scripted config changes without the TUI form.

**Work:**
- Add `config set` subcommand in `cmd/config.go`
- Validate key against allowed set (`provider`, `model`, `base`, `anthropic_api_key`, `openai_api_key`, `ollama_base_url`)
- Load existing config, update the key, write back via `config.Save()`

**Files:** `cmd/config.go`, `internal/config/save.go`

---

#### 5.3 `config init` — scaffold project-local config

**Goal:** One command to create a `.pr-pilot.toml` in the current directory.

**Work:**
- Add `config init` subcommand
- Write a commented template TOML with all keys and their defaults
- Abort if `.pr-pilot.toml` already exists (unless `--force`)

**Files:** `cmd/config.go`

---

## Suggested Build Order

| Step | Item | Effort |
|------|------|--------|
| 1 | 2.2 System/user message split | Small |
| 2 | 2.1 GitHub PR creation | Medium |
| 3 | 4.1 Diff truncation | Small |
| 4 | 3.2 Copy to clipboard | Small |
| 5 | 4.2 Configurable Ollama URL | Small |
| 6 | 3.1 Streaming + spinner | Medium |
| 7 | 5.1 Unit tests | Large |
| 8 | 5.2 `config set` | Small |
| 9 | 5.3 `config init` | Small |

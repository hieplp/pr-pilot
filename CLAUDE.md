# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o pr-pilot .

# Run without building
go run main.go [command]

# Tests
go test ./...
go test -v ./internal/provider/...   # single package
go test -race ./...                  # with race detector

# Lint (if golangci-lint is installed)
golangci-lint run
```

## Architecture

`pr-pilot` is a single-binary Go CLI. All business logic lives under `internal/`; `cmd/` is a thin layer that only parses flags and calls `internal/` functions.

### Data flow (commit command)

```
cmd/commit.go
  → internal/config.Load()          reads TOML / env vars / flag overrides
  → internal/git.StagedDiff()       shells out to `git diff --cached`
  → internal/prompt.CommitPrompt()  injects diff into system prompt
  → internal/provider.New()         selects Claude / OpenAI / Ollama
  → provider.Complete()             calls LLM API, returns text
  → fmt.Println()                   prints to stdout
```

The `pr` command follows the same flow but uses `git.BranchDiff(base)` + `git.CommitLog(base)` + `prompt.PRPrompt()`.

### Provider system

`internal/provider/provider.go` defines a two-method interface:

```go
type Provider interface {
    Complete(ctx context.Context, prompt string) (string, error)
    Name() string
}
```

`factory.go` maps the `--provider` flag string to a concrete implementation. To add a new provider: implement the interface in a new file, register it in `factory.go`'s `switch`.

### Config precedence

CLI flags → `PR_PILOT_*` env vars → `~/.config/pr-pilot/config.toml` (lowest).  
`config.Load()` returns a `*Config`; callers then call `cfg.Override(providerFlag, modelFlag)` to apply flag values on top.

### Key env vars

| Variable | Purpose |
|----------|---------|
| `ANTHROPIC_API_KEY` | Required for Claude (default provider) |
| `OPENAI_API_KEY` | Required for OpenAI |
| `PR_PILOT_PROVIDER` | Default provider (`claude`, `openai`, `ollama`) |
| `PR_PILOT_MODEL` | Default model name |

## Conventions

- `cmd/` files: flag parsing + error return only. No business logic.
- Git operations: always shell out via `os/exec` (see `internal/git/diff.go`), never import a git library.
- New subcommand: create `cmd/<name>.go`, define `cobra.Command`, register in `init()` via `rootCmd.AddCommand`.
- New config key: add default + env binding in `internal/config/config.go` using `viper.SetDefault` / `viper.BindEnv`.

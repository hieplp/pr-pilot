# Development Guide

## Prerequisites

- Go 1.24+
- Git
- An API key for at least one provider (Claude, OpenAI, or a running Ollama instance)

Verify your Go version:

```bash
go version
```

## Getting Started

```bash
git clone https://github.com/hieplp/pr-pilot
cd pr-pilot
go mod download
go run main.go --help
```

## Project Structure

```
pr-pilot/
├── main.go                    # entrypoint
├── cmd/
│   ├── root.go                # root command, global --provider / --model flags
│   ├── commit.go              # `pr-pilot commit` subcommand
│   └── pr.go                  # `pr-pilot pr` subcommand
└── internal/
    ├── provider/
    │   ├── provider.go        # Provider interface
    │   ├── claude.go          # Anthropic Claude implementation
    │   ├── openai.go          # OpenAI implementation
    │   └── ollama.go          # Ollama (local) implementation
    ├── git/
    │   └── diff.go            # git diff / log helpers
    ├── prompt/
    │   └── prompt.go          # prompt building and template rendering
    ├── config/
    │   └── config.go          # Viper-based config loader
    └── store/
        └── store.go           # SQLite usage log
```

All business logic lives under `internal/`. The `cmd/` layer is thin — it parses flags, calls `internal/` functions, and handles errors.

## Adding a New Provider

1. Create `internal/provider/<name>.go`.
2. Implement the `Provider` interface:

```go
type Provider interface {
    Complete(ctx context.Context, prompt string) (string, error)
    Name() string
}
```

3. Register it in the provider factory (the function that maps the `--provider` flag value to a concrete implementation).

Example skeleton:

```go
package provider

import "context"

type MyProvider struct {
    model string
}

func NewMyProvider(model string) *MyProvider {
    return &MyProvider{model: model}
}

func (p *MyProvider) Name() string { return "myprovider" }

func (p *MyProvider) Complete(ctx context.Context, prompt string) (string, error) {
    // call the API, return the text response
}
```

## Adding a New Subcommand

1. Create `cmd/<name>.go`.
2. Define a `cobra.Command` and register it in `init()`:

```go
package cmd

import "github.com/spf13/cobra"

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "One-line description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

## Configuration

The app reads config from (in order of precedence):

1. CLI flags (`--provider`, `--model`)
2. Environment variables (`PR_PILOT_PROVIDER`, `PR_PILOT_MODEL`)
3. Config file at `~/.config/pr-pilot/config.toml`

Example `config.toml`:

```toml
provider = "claude"
model    = ""
```

When adding a new config key, define it in `internal/config/config.go` and bind it with Viper:

```go
viper.SetDefault("my_key", "default_value")
viper.BindEnv("my_key", "PR_PILOT_MY_KEY")
```

## Running & Building

```bash
# Run without building
go run main.go commit

# Build binary
go build -o pr-pilot .

# Install to $GOPATH/bin
go install .
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/provider/...

# Run with verbose output
go test -v ./...

# Run with race detector
go test -race ./...
```

Mock the LLM in tests by implementing the `Provider` interface with a stub:

```go
type stubProvider struct{ response string }

func (s *stubProvider) Name() string { return "stub" }
func (s *stubProvider) Complete(_ context.Context, _ string) (string, error) {
    return s.response, nil
}
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | API key for Claude |
| `OPENAI_API_KEY` | API key for OpenAI |
| `PR_PILOT_PROVIDER` | Default provider (`claude`, `openai`, `ollama`) |
| `PR_PILOT_MODEL` | Default model name |

## Common Workflows

### Implementing `commit` end-to-end

1. In `internal/git/diff.go`, shell out to `git diff --cached` and return the output.
2. In `internal/prompt/prompt.go`, build the prompt string from the diff.
3. In `cmd/commit.go`, wire: read flags → load provider → get diff → build prompt → call `provider.Complete` → print result.

### Implementing `pr` end-to-end

Same as above but use `git diff <base>...HEAD` instead of `git diff --cached`.

## Code Conventions

- Keep `cmd/` thin — no business logic, only flag parsing and output.
- Return errors up the call stack; let Cobra print them.
- Shell out to `git` via `os/exec` rather than using a git library.
- Write a test for every new `internal/` package.

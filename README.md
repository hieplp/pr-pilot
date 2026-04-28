# pr-pilot

A CLI tool that generates commit messages and PR descriptions from git diffs using LLMs (Claude, OpenAI, Ollama).

## Installation

```bash
go install github.com/hieplp/pr-pilot@latest
```

Or build from source:

```bash
git clone https://github.com/hieplp/pr-pilot
cd pr-pilot
go build -o pr-pilot .
```

## Usage

```bash
# Generate a commit message from staged changes
pr-pilot commit

# Generate a PR description against main
pr-pilot pr

# Diff against a custom base branch
pr-pilot pr --base develop

# Use a specific provider and model
pr-pilot --provider openai --model gpt-4o commit
```

## Configuration

Set your API key via environment variable:

```bash
export ANTHROPIC_API_KEY=your-key   # for Claude (default)
export OPENAI_API_KEY=your-key      # for OpenAI
```

Or create a config file at `~/.config/pr-pilot/config.toml`:

```toml
provider = "claude"
model    = ""        # leave empty to use the provider default
```

## Supported Providers

| Flag value | Backend |
|------------|---------|
| `claude` (default) | Anthropic Claude via `anthropic-sdk-go` |
| `openai` | OpenAI via `go-openai` |
| `ollama` | Local Ollama (OpenAI-compatible endpoint) |

## Project Layout

```
pr-pilot/
├── main.go
├── cmd/
│   ├── root.go       # global --provider / --model flags
│   ├── commit.go     # staged diff → commit message
│   └── pr.go         # branch diff → PR description
└── internal/
    ├── provider/     # LLM provider abstraction + implementations
    ├── git/          # git diff / log helpers
    ├── prompt/       # prompt building and templates
    ├── config/       # Viper-based config loader
    └── store/        # SQLite usage log
```

## Development

```bash
# Run without building
go run main.go --help

# Run tests
go test ./...

# Build
go build -o pr-pilot .
```

## Roadmap

- [x] CLI skeleton (Cobra), `commit` and `pr` subcommands
- [x] `Provider` interface
- [ ] Claude / OpenAI / Ollama provider implementations
- [ ] Git diff and log helpers
- [ ] Config file + env var support
- [ ] Interactive TUI review before committing
- [ ] Pre-commit hook installer (`pr-pilot hook install`)
- [ ] Direct GitHub/GitLab PR creation
- [ ] Changelog generation
- [ ] Token cost dry-run (`--dry-run`)

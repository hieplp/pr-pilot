# Tech Stack Suggestions for pr-pilot

---

## Recommended: Go

The best fit for a CLI tool that ships as a single binary with no runtime dependency.

| Concern | Choice | Reason |
|---------|--------|--------|
| Language | **Go 1.22+** | Single static binary, fast startup (<50 ms), easy cross-platform builds |
| CLI framework | **[Cobra](https://github.com/spf13/cobra)** | Industry standard for Go CLIs; subcommands, flags, completions built-in |
| Config | **[Viper](https://github.com/spf13/viper)** | TOML/YAML/env var config with zero boilerplate; pairs with Cobra |
| LLM — Claude | **[anthropic-sdk-go](https://github.com/anthropic-ai/anthropic-sdk-go)** | Official Go SDK, streaming support, prompt caching |
| LLM — OpenAI | **[go-openai](https://github.com/sashabaranov/go-openai)** | Mature, covers OpenAI + any compatible endpoint (Ollama, Groq, etc.) |
| Git operations | **`os/exec` + git CLI** | Shell out to git for diff/log; avoid a git library to stay simple |
| Interactive TUI | **[bubbletea](https://github.com/charmbracelet/bubbletea)** | Pick-list for candidates, editor preview, spinner while waiting for LLM |
| Terminal styling | **[lipgloss](https://github.com/charmbracelet/lipgloss)** | Colored output, boxes, tables — same ecosystem as bubbletea |
| Local storage | **[modernc/sqlite](https://gitlab.com/cznic/sqlite)** | Pure-Go SQLite (no cgo), for usage stats/history |
| Secret detection | **[gitleaks](https://github.com/gitleaks/gitleaks)** (as a lib) | Reuse its regex ruleset to scrub diff before sending to LLM |
| Testing | **`testing` + [testify](https://github.com/stretchr/testify)** | Standard + assertions; mock LLM with an httptest server |
| Distribution | **[goreleaser](https://goreleaser.com/)** | Builds for Linux/macOS/Windows, publishes to GitHub Releases + Homebrew |

### Project layout

```
pr-pilot/
├── cmd/
│   ├── root.go          # cobra root + global flags
│   ├── commit.go        # pr-pilot commit
│   ├── pr.go            # pr-pilot pr
│   └── hook.go          # pr-pilot hook install/uninstall
├── internal/
│   ├── provider/        # LLM provider abstraction
│   │   ├── provider.go  # interface
│   │   ├── claude.go
│   │   ├── openai.go
│   │   └── ollama.go
│   ├── git/             # git diff, log, branch helpers
│   ├── prompt/          # prompt building + template rendering
│   ├── config/          # viper-based config loader
│   └── store/           # sqlite usage log
├── .goreleaser.yaml
└── main.go
```

---

## Alternative: Python

Better if you want to iterate fast on prompts and LLM logic, or already know Python well.

| Concern | Choice |
|---------|--------|
| CLI framework | **[Typer](https://typer.tiangolo.com/)** (built on Click) |
| LLM | **[anthropic](https://pypi.org/project/anthropic/)** + **[openai](https://pypi.org/project/openai/)** official SDKs |
| Config | **[dynaconf](https://www.dynaconf.com/)** or plain `tomllib` (stdlib 3.11+) |
| TUI / prompts | **[rich](https://github.com/Textualize/rich)** + **[questionary](https://github.com/tmbo/questionary)** |
| Distribution | **[uv](https://github.com/astral-sh/uv)** for packaging, **[PyInstaller](https://pyinstaller.org/)** for single binary |

Downside: requires Python runtime on user machines; `uv tool install` mitigates this.

---

## Alternative: TypeScript (Node)

Good if you want to share code with a future web UI or VS Code extension.

| Concern | Choice |
|---------|--------|
| Runtime | **[Bun](https://bun.sh/)** — fast startup, single-file compile |
| CLI framework | **[commander](https://github.com/tj/commander.js/)** or **[oclif](https://oclif.io/)** |
| LLM | **[@anthropic-ai/sdk](https://www.npmjs.com/package/@anthropic-ai/sdk)** + **[openai](https://www.npmjs.com/package/openai)** |
| TUI | **[ink](https://github.com/vadimdemedes/ink)** (React for terminal) |
| Distribution | **[pkg](https://github.com/vercel/pkg)** or Bun's `bun build --compile` |

Downside: heavier binary, slower cold start compared to Go.

---

## Decision Guide

```
Do you want a single binary with no install friction?  → Go
Do you want to move fast on prompt logic and don't mind pip/uv? → Python
Do you plan to share code with a VS Code extension or web UI? → TypeScript
```

**Recommendation: Go.** A developer tool that wraps git should feel instant, ship as one file, and be easy to install via Homebrew or `curl | sh`. Go hits all three.

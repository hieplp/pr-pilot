# Feature Suggestions for pr-pilot

A CLI that generates commit messages and PR descriptions from git diffs using LLMs (Claude, OpenAI, local models, etc.).

---

## Core Generation

### Commit Message Generation
- Analyze staged diff (`git diff --cached`) and generate a Conventional Commits-style message
- Support scopes, breaking change footers, and multi-line bodies automatically
- Option to generate multiple candidates and let the user pick

### PR Description Generation
- Diff against a base branch (default: `main`/`master`) and produce a structured PR body:
  - Summary, motivation, changes breakdown, test plan, screenshots placeholder
- Detect issue/ticket references from branch name (`feat/PROJ-123-...`) and auto-link them

### Changelog Entry Generation
- From a range of commits or a tag-to-HEAD diff, produce a `CHANGELOG.md` section grouped by type (feat, fix, chore, etc.)

---

## LLM Provider Support

### Provider Abstraction
- Pluggable provider system: Claude (Anthropic), OpenAI, Gemini, Ollama (local), any OpenAI-compatible endpoint
- Per-project or global config file (`~/.config/pr-pilot/config.toml`) to set default provider and model
- `--provider` and `--model` flags to override per invocation

### Cost & Token Awareness
- Show estimated token count and cost before sending (with `--dry-run`)
- Warn when diff is large; offer to truncate or summarize file-by-file before sending

### Prompt Caching
- Use Anthropic prompt caching for repeated system prompts to reduce latency and cost on Claude

---

## Git Integration

### Staged-Only vs Full Branch Diff
- `commit` subcommand: works on staged changes
- `pr` subcommand: diffs current branch against base branch
- `review` subcommand: summarizes each commit in a branch individually

### Pre-commit Hook Mode
- Install as a `prepare-commit-msg` git hook so the message is pre-filled whenever you run `git commit`
- `pr-pilot hook install / uninstall`

### Branch Name Suggestions
- From a description or staged diff, suggest a git branch name following a project convention (e.g., `feat/short-slug`)

---

## Output & Workflow

### Interactive Review
- Show generated output in `$EDITOR` or a TUI diff view before committing/pushing
- Accept / edit / regenerate options without leaving the terminal

### Direct Push & PR Creation
- After generating a PR description, optionally open a GitHub/GitLab/Bitbucket PR with one command
- Use provider APIs (gh CLI, GitLab API) — no vendor lock-in

### Output Formats
- Plain text (default), JSON (for scripting), Markdown, clipboard copy (`--copy`)
- Template support: users supply a Handlebars/Jinja template to control the exact output shape

---

## Configuration & Customization

### Project-level Config
- `.pr-pilot.toml` in repo root to set base branch, provider, commit style, template paths
- Override any setting with env vars (`PR_PILOT_PROVIDER=ollama`)

### Commit Convention Presets
- Built-in presets: Conventional Commits, Angular, gitmoji, semantic
- Custom regex patterns for teams with house style

### Context Injection
- `--context <file>` to attach extra context (e.g., a JIRA ticket body, a design doc snippet) that the LLM uses when writing the description
- Auto-read `CONTRIBUTING.md` or `PR_TEMPLATE.md` from the repo and include in the prompt

---

## Quality & Safety

### Diff Filtering
- Exclude generated files, lock files, and vendor directories from the diff before sending to the LLM (configurable ignore list)
- Strip secrets/tokens from diff using a regex allowlist before sending

### Hallucination Guard
- Verify that file names and symbols mentioned in the generated message actually appear in the diff
- Warn (or strip) any line that references a file not in the changeset

### Retry & Fallback
- Automatic retry with exponential backoff on rate-limit errors
- Fallback chain: if primary provider fails, try secondary (e.g., Claude → OpenAI → local Ollama)

---

## Developer Experience

### Shell Completions
- Generate completions for bash, zsh, fish (`pr-pilot completions zsh`)

### Verbose / Debug Mode
- `--verbose` prints the exact prompt sent and raw response received
- `--debug` shows token counts, latency, and provider metadata

### Stats & History
- Local SQLite log of generations: timestamp, provider, tokens, cost, diff hash
- `pr-pilot stats` shows usage summary per provider/project

---

## Roadmap Priorities (suggested order)

| Priority | Feature |
|----------|---------|
| 1 | Commit message generation (staged diff → Conventional Commits) |
| 2 | Multi-provider abstraction (Claude + OpenAI + Ollama) |
| 3 | PR description generation (branch diff → structured body) |
| 4 | Config file + env var support |
| 5 | Pre-commit hook installer |
| 6 | Interactive editor review before committing |
| 7 | Direct GitHub/GitLab PR creation |
| 8 | Changelog generation |
| 9 | Cost/token dry-run estimation |
| 10 | Local usage stats |

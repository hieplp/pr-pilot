package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Action represents what the user chose in the review screen.
type Action int

const (
	ActionAccept Action = iota
	ActionEdit
	ActionRegenerate
	ActionQuit
)

// Result is returned by Review.
type Result struct {
	Action  Action
	Content string
}

// Review shows generated content in a TUI and waits for the user to accept,
// edit, regenerate, or quit. On ActionEdit the external $EDITOR is opened and
// ActionAccept is returned with the edited content.
func Review(content string) (Result, error) {
	p := tea.NewProgram(reviewModel{content: content}, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return Result{}, err
	}

	m := final.(reviewModel)
	result := Result{Action: m.action, Content: content}

	if m.action == ActionEdit {
		edited, err := openInEditor(content)
		if err != nil {
			return Result{}, err
		}
		result.Content = edited
		result.Action = ActionAccept
	}

	return result, nil
}

// ── bubbletea model ──────────────────────────────────────────────────────────

type reviewModel struct {
	content string
	action  Action
	done    bool
}

func (m reviewModel) Init() tea.Cmd { return nil }

func (m reviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
	case "a", "enter":
		m.action = ActionAccept
	case "e":
		m.action = ActionEdit
	case "r":
		m.action = ActionRegenerate
	case "q", "ctrl+c", "esc":
		m.action = ActionQuit
	default:
		return m, nil
	}
	m.done = true
	return m, tea.Quit
}

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)
)

func (m reviewModel) View() string {
	if m.done {
		return ""
	}
	hints := fmt.Sprintf(
		"%s accept  %s edit  %s regenerate  %s quit",
		keyStyle.Render("[a/↵]"),
		keyStyle.Render("[e]"),
		keyStyle.Render("[r]"),
		keyStyle.Render("[q]"),
	)
	return labelStyle.Render("Generated output") + "\n\n" +
		boxStyle.Render(m.content) + "\n\n" +
		hintStyle.Render(hints) + "\n"
}

// ── editor helper ────────────────────────────────────────────────────────────

func openInEditor(content string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	f, err := os.CreateTemp("", "pr-pilot-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return "", err
	}
	f.Close()

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	edited, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(edited)), nil
}

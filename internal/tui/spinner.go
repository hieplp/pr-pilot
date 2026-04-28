package tui

import (
	"errors"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type workDoneMsg struct {
	content string
	err     error
}

type spinModel struct {
	spinner spinner.Model
	label   string
	done    *workDoneMsg
}

func (m spinModel) Init() tea.Cmd { return m.spinner.Tick }

func (m spinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case workDoneMsg:
		m.done = &v
		return m, tea.Quit
	case tea.KeyMsg:
		if v.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

var spinLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

func (m spinModel) View() string {
	if m.done != nil {
		return ""
	}
	return m.spinner.View() + " " + spinLabelStyle.Render(m.label) + "\n"
}

// Spin shows an animated spinner with label while work runs in the background.
// Returns when work completes or ctrl+c is pressed.
func Spin(label string, work func() (string, error)) (string, error) {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := spinModel{spinner: s, label: label}
	p := tea.NewProgram(m)

	go func() {
		content, err := work()
		p.Send(workDoneMsg{content: content, err: err})
	}()

	final, err := p.Run()
	if err != nil {
		return "", err
	}

	sm := final.(spinModel)
	if sm.done == nil {
		return "", errors.New("cancelled")
	}
	return sm.done.content, sm.done.err
}

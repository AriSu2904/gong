package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type taskResultMsg struct{ err error }

type model struct {
	spinner  spinner.Model
	taskFunc func() error
	status   string
	err      error
}

func newModel(taskFunc func() error, status string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return model{
		spinner:  s,
		taskFunc: taskFunc,
		status:   status,
	}
}

func RunSpinner(statusText string, task func() error) error {
	initialModel := newModel(task, statusText)
	p := tea.NewProgram(initialModel)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("gagal menjalankan TUI: %w", err)
	}

	if m, ok := finalModel.(model); ok && m.err != nil {
		return m.err
	}

	return nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.runTask)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case taskResultMsg:
		m.err = msg.err
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.status)
}

func (m model) runTask() tea.Msg {
	err := m.taskFunc()
	return taskResultMsg{err: err}
}

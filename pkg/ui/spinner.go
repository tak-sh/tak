package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &Spinner{}

func NewSpinner() *Spinner {
	out := &Spinner{
		spinner: spinner.New(),
	}

	out.spinner.Spinner = spinner.Dot

	return out
}

type Spinner struct {
	spinner spinner.Model
}

func (s *Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

func (s *Spinner) View() string {
	return s.spinner.View()
}

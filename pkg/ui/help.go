package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tak-sh/tak/pkg/ui/keyregistry"
	"strings"
)

func newHelpModel() *HelpModel {
	return &HelpModel{
		Help: help.New(),
		Keys: keyregistry.DefaultKeys,
	}
}

var _ tea.Model = &HelpModel{}

type HelpModel struct {
	Help help.Model
	Keys help.KeyMap

	width  int
	height int
}

func (h *HelpModel) Init() tea.Cmd {
	return nil
}

func (h *HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.Help.Width = msg.Width
		h.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyregistry.DefaultKeys.Help):
			h.Help.ShowAll = !h.Help.ShowAll
		case key.Matches(msg, keyregistry.DefaultKeys.Quit):
			return h, tea.Quit
		}
	}

	return h, nil
}

func (h *HelpModel) View() string {
	helpView := h.Help.View(h.Keys)
	height := 8 - strings.Count(helpView, "\n")

	return "\n" + strings.Repeat("\n", height) + helpView
}

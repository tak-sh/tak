package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tak-sh/tak/pkg/ui/keyregistry"
	"strings"
)

func newHelpModel() *helpModel {
	return &helpModel{
		Help: help.New(),
		Keys: keyregistry.DefaultKeys,
	}
}

var _ tea.Model = &helpModel{}

type helpModel struct {
	Help help.Model
	Keys *keyregistry.KeyMap

	width  int
	height int
}

func (h *helpModel) Init() tea.Cmd {
	return nil
}

func (h *helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.Help.Width = msg.Width
		h.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.Keys.Help):
			h.Help.ShowAll = !h.Help.ShowAll
		case key.Matches(msg, h.Keys.Quit):
			return h, tea.Quit
		}
	}

	return h, nil
}

func (h *helpModel) View() string {
	helpView := h.Help.View(h.Keys)
	height := 8 - strings.Count(helpView, "\n")

	return "\n" + strings.Repeat("\n", height) + helpView
}

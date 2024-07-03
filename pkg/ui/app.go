package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tak-sh/tak/pkg/ui/keyregistry"
	"github.com/tak-sh/tak/pkg/utils/bubbleutils"
)

var _ tea.Model = &App{}

func NewApp(children ...tea.Model) *App {
	out := &App{
		Children: children,
		Help:     newHelpModel(),
	}

	return out
}

type App struct {
	Children     []tea.Model
	OnReady      chan bool
	Help         *HelpModel
	windowWidth  int
	windowHeight int
}

func (b *App) Init() tea.Cmd {
	return bubbleutils.InitAll(append(b.Children, b.Help))
}

func (b *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.WindowSizeMsg:
		b.windowHeight = t.Height
		b.windowWidth = t.Width
	case tea.KeyMsg:
		if key.Matches(t, keyregistry.DefaultKeys.Quit) {
			return b, tea.Quit
		}
	}

	cmd := bubbleutils.UpdateAll(msg, append(b.Children, b.Help))

	return b, cmd
}

func (b *App) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, append(bubbleutils.RenderAll(b.Children), b.Help.View())...)
}

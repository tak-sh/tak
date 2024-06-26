package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/component"
	"github.com/tak-sh/tak/pkg/ui/keyregistry"
	"github.com/tak-sh/tak/pkg/utils/bubbleutils"
)

var _ tea.Model = &App{}

type SubmitEvent struct {
	ID  string
	Val *v1beta1.Value
}

type OnSubmitFunc func(s *SubmitEvent)

func NewApp(onSubmit OnSubmitFunc, msg string) *App {
	out := &App{
		Children:               []tea.Model{},
		OnSubmit:               onSubmit,
		Spinner:                NewSpinner(),
		DefaultProgressMessage: msg,
		ProgressMessage:        msg,
		help:                   newHelpModel(),
	}

	out.Spinner.spinner.Style = SpinnerStyle

	return out
}

type App struct {
	Children               []tea.Model
	OnSubmit               OnSubmitFunc
	Spinner                *Spinner
	ProgressMessage        string
	DefaultProgressMessage string
	OnReady                chan bool

	showSpinner  bool
	help         *helpModel
	windowWidth  int
	windowHeight int
}

func (b *App) Init() tea.Cmd {
	return bubbleutils.InitAll(append(b.Children, b.Spinner, b.help))
}

func (b *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	defer func() {
		b.showSpinner = len(b.Children) == 0
	}()
	switch t := msg.(type) {
	case tea.WindowSizeMsg:
		b.windowHeight = t.Height
		b.windowWidth = t.Width
	case tea.KeyMsg:
		if key.Matches(t, keyregistry.DefaultKeys.Quit) {
			return b, tea.Quit
		} else if key.Matches(t, keyregistry.DefaultKeys.Submit) {
			cmd := bubbleutils.UpdateAll(component.SyncStateMsg(func(id string, v *v1beta1.Value) {
				b.OnSubmit(&SubmitEvent{
					ID:  id,
					Val: v,
				})
			}), b.Children)
			b.Children = nil
			b.ProgressMessage = b.DefaultProgressMessage
			return b, cmd
		}
	case SetChildrenMsg:
		b.Children = t.Models
		return b, bubbleutils.InitAll(b.Children)
	case AppendChildrenMsg:
		b.Children = append(b.Children, t.Models...)
		return b, bubbleutils.InitAll(b.Children)
	case UpdateProgressMsg:
		b.ProgressMessage = t.Msg
	}

	cmd := bubbleutils.UpdateAll(msg, append(b.Children, b.help, b.Spinner))

	return b, cmd
}

func (b *App) View() string {
	body := make([]string, 0, len(b.Children)+1)

	if b.showSpinner {
		body = append(body, lipgloss.JoinHorizontal(lipgloss.Left, b.Spinner.View(), ProgressMessageStyle.Render(b.ProgressMessage)))
	} else {
		body = append(body, lipgloss.JoinVertical(lipgloss.Left, bubbleutils.RenderAll(b.Children)...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, append(body, b.help.View())...)
}

type SetChildrenMsg struct {
	Models []tea.Model
}

type AppendChildrenMsg struct {
	Models []tea.Model
}

type UpdateProgressMsg struct {
	Msg string
}

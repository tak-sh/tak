package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tak-sh/tak/pkg/headless/component"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/utils/stringutils"
	"log/slog"
)

type OnRenderFunc func(id string)

func NewScriptComponent(accountName string, str renderer.Stream, eq engine.EventQueue, logger *slog.Logger) *ScriptComponent {
	out := &ScriptComponent{
		Stream:                 str,
		ScriptEvents:           eq,
		Logger:                 logger,
		Spinner:                NewSpinner(),
		DefaultProgressMessage: fmt.Sprintf("Adding your %s account...", accountName),
		showSpinner:            true,
	}

	out.ProgressMessage = out.DefaultProgressMessage
	out.Spinner.spinner.Style = SpinnerStyle

	return out
}

var _ tea.Model = &ScriptComponent{}

type ScriptComponent struct {
	Stream       renderer.Stream
	ScriptEvents engine.EventQueue
	OnRenderFunc OnRenderFunc
	Logger       *slog.Logger

	// visual components
	Child                  tea.Model
	Spinner                *Spinner
	ProgressMessage        string
	DefaultProgressMessage string
	showSpinner            bool
}

func (s *ScriptComponent) Init() tea.Cmd {
	cmds := []tea.Cmd{
		s.waitEventQueueMsg(),
		s.waitRenderQueueMsg(),
		s.Spinner.Init(),
	}

	if s.Child != nil {
		cmds = append(cmds, s.Child.Init())
	}

	return tea.Batch(cmds...)
}

func (s *ScriptComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	defer func() {
		s.showSpinner = s.Child == nil
	}()

	cmds := make([]tea.Cmd, 0, 2)
	switch t := msg.(type) {
	case component.OnSubmitMsg:
		s.Child = nil
		s.ProgressMessage = s.DefaultProgressMessage
		s.Stream.Respond(&renderer.Response{
			ID:    t.Id,
			Value: t.Value,
		})
	case RenderModelMsg:
		s.Child = t.Model
		cmds = append(cmds, s.Child.Init(), s.onRender(t.Model), s.waitRenderQueueMsg())
	case OnScriptEventMsg:
		switch m := t.Event.(type) {
		case *engine.NextInstructionEvent:
			if _, ok := m.Instruction.(*step.PromptAction); !ok {
				s.ProgressMessage = stringutils.Capitalize(m.String())
				s.Child = nil
			}
		}

		cmds = append(cmds, s.waitEventQueueMsg())
	}

	if s.Child != nil {
		var cmd tea.Cmd
		s.Child, cmd = s.Child.Update(msg)
		cmds = append(cmds, cmd)
	}

	_, spinnerCmd := s.Spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return s, tea.Batch(cmds...)
}

func (s *ScriptComponent) View() string {
	if s.showSpinner {
		return lipgloss.JoinHorizontal(lipgloss.Left, s.Spinner.View(), ProgressMessageStyle.Render(s.ProgressMessage))
	} else if s.Child != nil {
		return s.Child.View()
	}
	return ""
}

func (s *ScriptComponent) onRender(r renderer.Model) tea.Cmd {
	if s.OnRenderFunc == nil {
		return nil
	}

	return func() tea.Msg {
		s.OnRenderFunc(r.GetId())
		return nil
	}
}

func (s *ScriptComponent) waitRenderQueueMsg() tea.Cmd {
	return func() tea.Msg {
		r, ok := <-s.Stream.RenderQueue()
		if !ok {
			s.Logger.Info("Render queue closed.")
			return tea.Quit()
		}
		return RenderModelMsg{Model: r}
	}
}

func (s *ScriptComponent) waitEventQueueMsg() tea.Cmd {
	return func() tea.Msg {
		e, ok := <-s.ScriptEvents
		if !ok {
			s.Logger.Info("Event queue has closed.")
			return tea.Quit()
		}
		return OnScriptEventMsg{Event: e}
	}
}

type RenderModelMsg struct {
	Model renderer.Model
}

func (r RenderModelMsg) String() string {
	return fmt.Sprintf("display %s component", r.Model.GetId())
}

type OnScriptEventMsg struct {
	Event engine.Event
}

func (o OnScriptEventMsg) String() string {
	return o.Event.String()
}

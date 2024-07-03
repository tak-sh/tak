package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/tak-sh/tak/pkg/debug"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/protoenc"
	"github.com/tak-sh/tak/pkg/ui/keyregistry"
	"slices"
	"strings"
	"time"
)

var _ tea.Model = &DebugComponent{}

func NewDebugComponent(s debug.Stepper, sc *ScriptComponent) *DebugComponent {
	out := &DebugComponent{
		SC:            sc,
		Debug:         s,
		DebugView:     &DebugView{},
		UpdateHistory: NewUpdateHistory(10),
	}

	return out
}

type DebugComponent struct {
	SC            *ScriptComponent
	Debug         debug.Stepper
	DebugView     *DebugView
	UpdateHistory *UpdateHistory
}

func (d *DebugComponent) Init() tea.Cmd {
	return d.SC.Init()
}

func (d *DebugComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 1)
	switch t := msg.(type) {
	case tea.WindowSizeMsg:
		d.DebugView.Width = t.Width
		d.DebugView.Height = t.Height
	case OnScriptEventMsg:
		switch t := t.Event.(type) {
		case *engine.NextInstructionEvent:
			if v, ok := t.Instruction.(*step.Step); ok {
				d.DebugView.setStep(v)
			}
			d.DebugView.setStepData(t.Context.TemplateData)
		}
	case tea.KeyMsg:
		if key.Matches(t, keyregistry.DebugKeys.Next, keyregistry.DefaultKeys.Submit) {
			d.Debug.Step()
		} else if key.Matches(t, keyregistry.DebugKeys.Previous) {
			d.Debug.PreviousStep()
		} else if key.Matches(t, keyregistry.DebugKeys.Reload) {
			d.Debug.Replay()
		}
	}
	mod, cmd := d.SC.Update(msg)
	cmds = append(cmds, cmd)
	d.SC = mod.(*ScriptComponent)

	if d.UpdateHistory != nil {
		mod, cmd = d.UpdateHistory.Update(msg)
		cmds = append(cmds, cmd)
		d.UpdateHistory = mod.(*UpdateHistory)
	}

	return nil, tea.Batch(cmds...)
}

func (d *DebugComponent) View() string {
	t := table.New().Width(d.DebugView.Width).Height(d.DebugView.Height)
	rows := []string{
		d.SC.View(), d.DebugView.String(),
	}
	if d.UpdateHistory != nil {
		rows = append(rows, d.UpdateHistory.View())
	}
	t.Row(rows...)
	t.Border(lipgloss.HiddenBorder())
	return t.Render()
}

var (
	StepStyle            = lipgloss.NewStyle().Foreground(ComplementaryColor)
	DebugViewBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), false, false, false, true).BorderForeground(SecondaryColor)
)

var _ fmt.Stringer = &DebugView{}

type DebugView struct {
	Step            *step.Step
	StepData        *engine.TemplateData
	Width           int
	Height          int
	stepContent     string
	stepDataContent string
}

func (d *DebugView) setStep(s *step.Step) {
	d.Step = s
	b, _ := protoenc.MarshalYAML(s)
	d.stepContent = string(b)
}

func (d *DebugView) setStepData(data *engine.TemplateData) {
	d.StepData = data
	strs := make([]string, 0, len(data.Step))
	for k, v := range data.Step {
		strs = append(strs, fmt.Sprintf("step.%s = %s", k, v))
	}
	slices.Sort(strs)

	d.stepDataContent = strings.Join(strs, "\n")
}

func (d *DebugView) String() string {
	stepStr := d.stepContent
	if stepStr == "" {
		stepStr = fmt.Sprintf("Press %s to run the next step", strings.Join(keyregistry.DebugKeys.Next.Keys(), "or"))
	}

	out := strings.Join([]string{
		"Current step:",
		StepStyle.Render(Indent(stepStr, 2)),
		"Available values:",
		StepStyle.Render(Indent(d.stepDataContent, 2)),
	}, "\n")

	return DebugViewBorderStyle.Render(out)
}

var _ tea.Model = &UpdateHistory{}

func NewUpdateHistory(limit int) *UpdateHistory {
	out := &UpdateHistory{
		Limit:         limit,
		UpdateHistory: make([]*UpdateEntry, 0, limit),
		MessageFilter: DefaultUpdateFilter(),
		Style:         lipgloss.NewStyle().Foreground(SecondaryColor).Border(lipgloss.RoundedBorder(), false, false, false, true),
	}

	return out
}

func DefaultUpdateFilter() MessageFilter {
	return func(msg tea.Msg) bool {
		switch msg.(type) {
		case OnScriptEventMsg, tea.KeyMsg, RenderModelMsg:
			return true
		}
		return false
	}
}

type MessageFilter func(msg tea.Msg) bool

type UpdateHistory struct {
	Limit         int
	UpdateHistory []*UpdateEntry
	MessageFilter MessageFilter
	Style         lipgloss.Style
}

type UpdateEntry struct {
	Time time.Time
	Msg  fmt.Stringer
}

func (e *UpdateEntry) String() string {
	return fmt.Sprintf("%s - %s", e.Time.Format("3:04 05.000"), e.Msg.String())
}

func (u *UpdateHistory) Init() tea.Cmd { return nil }

func (u *UpdateHistory) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if u.MessageFilter == nil || !u.MessageFilter(msg) {
		return u, nil
	}

	entry := &UpdateEntry{
		Time: time.Now(),
		Msg:  msg.(fmt.Stringer),
	}
	if len(u.UpdateHistory) == u.Limit {
		for i := 0; i < u.Limit-1; i++ {
			u.UpdateHistory[i] = u.UpdateHistory[i+1]
		}
		u.UpdateHistory[u.Limit-1] = entry
	} else {
		u.UpdateHistory = append(u.UpdateHistory, entry)
	}
	return u, nil
}

func (u *UpdateHistory) View() string {
	strs := make([]string, 0, len(u.UpdateHistory))
	for _, v := range u.UpdateHistory {
		strs = append(strs, v.String())
	}

	return u.Style.Render(lipgloss.JoinVertical(lipgloss.Left, strs...))
}

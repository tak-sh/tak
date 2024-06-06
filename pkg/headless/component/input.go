package component

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/utils/ptr"
)

func NewInput(i *v1beta1.Component_Input) *Input {
	out := &Input{
		Component_Input: i,
	}

	return out
}

var _ Component = &Input{}

type Input struct {
	*v1beta1.Component_Input
}

func (i *Input) Render(_ *engine.Context, props *Props) renderer.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &InputModel{
		Props: props,
		Text:  ti,
		Comp:  i.Component_Input,
	}
}

func (i *Input) ToProto() *v1beta1.Component {
	return &v1beta1.Component{Input: i.Component_Input}
}

func (i *Input) Validate() error {
	return nil
}

var _ renderer.Model = &InputModel{}

type InputModel struct {
	Props *Props
	Text  textinput.Model
	Comp  *v1beta1.Component_Input
}

func (i *InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (i *InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch t := msg.(type) {
	case SyncStateMsg:
		t(i.Props.ID, &v1beta1.Value{Str: ptr.Ptr(i.Text.Value())})
	default:
		i.Text, cmd = i.Text.Update(msg)
	}

	return i, cmd
}

func (i *InputModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, i.Props.Title, i.Props.Description, i.Text.View())
}

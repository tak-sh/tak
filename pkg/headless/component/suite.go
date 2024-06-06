package component

import (
	"bytes"
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/bubbleutils"
	"github.com/tak-sh/tak/pkg/utils/testutils"
	"io"
	"strings"
	"sync/atomic"
)

type TestSuite struct {
	suite.Suite
}

func (s *TestSuite) NewApp(store *engine.TemplateData, props *Props, components ...Component) *Commander {
	models := make([]tea.Model, 0, len(components))
	c, _ := engine.NewContext(context.Background(), nil, engine.ContextOpts{})
	c.TemplateData = c.TemplateData.Merge(store)

	for _, v := range components {
		models = append(models, v.Render(c, props))
	}

	onRender := make(chan []string, 1)
	app := &TestApp{
		Components: models,
		Renders:    make([]string, len(models)),
		OnRender: func(renders []string) {
			onRender <- renders
		},
	}

	output := bytes.NewBuffer([]byte{})
	input := bytes.NewBuffer([]byte{})
	prog := tea.NewProgram(app, tea.WithOutput(output), tea.WithInput(input))

	commander := &Commander{
		program:      prog,
		app:          app,
		onRenderChan: onRender,
		suite:        s,
		input:        input,
		output:       output,
	}

	return commander
}

func (s *TestSuite) EqualDropdownItems(expected []*dropdownItem, actual []list.Item, args ...any) bool {
	act := make([]*dropdownItem, len(actual))
	for i, v := range actual {
		act[i] = v.(*dropdownItem)
	}

	expComp := make([]*v1beta1.Component_Dropdown_Option, len(expected))
	actComp := make([]*v1beta1.Component_Dropdown_Option, len(act))
	for i, v := range expected {
		expComp[i] = v.comp
	}
	for i, v := range act {
		actComp[i] = v.comp
	}

	return testutils.AllEmpty(&s.Suite, testutils.EqualProtos(expComp, actComp), args...)
}

type Output interface {
	fmt.Stringer
	io.Writer
}

type Commander struct {
	suite        *TestSuite
	program      *tea.Program
	app          *TestApp
	onRenderChan chan []string
	output       Output
	input        *bytes.Buffer
}

func (c *Commander) Start() {
	go func() {
		defer func() {
			close(c.onRenderChan)
		}()
		_, err := c.program.Run()
		if err != nil {
			c.suite.FailNow(err.Error())
		}
	}()
}

func (c *Commander) Renders() chan []string {
	return c.onRenderChan
}

func (c *Commander) Stop() {
	c.program.Quit()
}

func (c *Commander) EqualOutput(expected string) bool {
	out := c.output.String()
	return c.suite.Equal(expected, out)
}

var _ tea.Model = &TestApp{}

type TestApp struct {
	Components  []tea.Model
	Renders     []string
	RenderCount atomic.Int32
	OnRender    func(renders []string)
}

func (t *TestApp) Init() tea.Cmd {
	return bubbleutils.InitAll(t.Components)
}

func (t *TestApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := bubbleutils.UpdateAll(msg, t.Components)
	return t, cmd
}

func (t *TestApp) View() string {
	t.RenderCount.Add(1)
	renders := make([]string, 0, len(t.Components))
	for i, v := range t.Components {
		t.Renders[i] = v.View()
	}

	if t.OnRender != nil {
		t.OnRender(t.Renders)
	}

	return strings.Join(renders, "\n")
}

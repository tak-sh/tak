package script

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/headless/step/stepper"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/ui"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
)

type ScriptTestSuite struct {
	suite.Suite
}

func (s *ScriptTestSuite) TestRun() {
	ser := newHTMLServer()
	defer ser.Close()

	sc := &v1beta1.Script{
		Steps: []*v1beta1.Step{
			{Action: newNavAction(path.Join(ser.URL, "testdata"))},
			{Action: newMouseClick("a[href*='table.html']")},
			{Id: ptr.Ptr("selected"), Action: newPromptAction(&v1beta1.Prompt{
				Title: "fruit",
				Component: newFromDropdownComponent(newEachSelector("table > tbody > tr", "td:first-child"), &v1beta1.Component_Dropdown_Option{
					Value: "{{ element.data }}",
				}),
			})},
			{Id: ptr.Ptr("input"), Action: newInput("input[id='test_input1']", "{{ step.selected }}")},
			{Action: newMouseClick("a[href*='done.html']")},
		},
		Signals: []*v1beta1.ConditionalSignal{
			{
				If:     "{{ '/done.html' in browser.url }}",
				Signal: v1beta1.ConditionalSignal_success,
			},
		},
	}

	r, w := io.Pipe()
	str := renderer.NewStream()
	eq := engine.NewEventQueue()
	bubble := ui.NewScriptComponent("derp", str, eq, slog.Default())
	bubble.OnRenderFunc = func(id string) {
		if id == "selected" {
			_, _ = w.Write([]byte(tea.KeyDown.String()))
			_, _ = w.Write([]byte(tea.KeyDown.String()))
			_, _ = w.Write([]byte(tea.KeyEnter.String()))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app := ui.NewApp(bubble)
	p := tea.NewProgram(
		app,
		tea.WithContext(ctx),
		tea.WithAltScreen(),
		tea.WithInput(r),
		tea.WithOutput(os.Stdout),
	)

	_, _ = ui.Run(ctx, p)
	c, _ := engine.NewContext(ctx, str, engine.NewEvaluator(eq, 1*time.Second), engine.ContextOpts{})
	comp, err := New(sc)
	if !s.NoError(err) {
		return
	}

	stper := stepper.New(comp.Signals, comp.Steps)

	doneCtx := RunAsync(c, comp, stper, WithPostRunFunc(func(c *engine.Context, st *step.Step) error {
		if st.GetId() == "input" {
			var val string
			_ = chromedp.Evaluate(`document.getElementById("test_input1").value`, &val).Do(c.Context)
			s.Equal("Orange", val)
		}
		return nil
	}))

	<-doneCtx.Done()

	s.EqualError(context.Canceled, context.Cause(doneCtx).Error())
	s.Equal("Orange", c.TemplateData.GetStepVal("selected"))
}

func (s *ScriptTestSuite) TestRunNoSuccessCondition() {
	ser := newHTMLServer()
	defer ser.Close()

	sc := &v1beta1.Script{
		Steps: []*v1beta1.Step{
			{Action: newNavAction(path.Join(ser.URL, "testdata"))},
		},
		Signals: []*v1beta1.ConditionalSignal{
			{
				If:     "{{ 'done.html' in browser.url }}",
				Signal: v1beta1.ConditionalSignal_success,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, _ := engine.NewContext(ctx, renderer.NewStream(), engine.NewEvaluator(engine.NewEventQueue(), 1*time.Second), engine.ContextOpts{})
	comp, err := New(sc)
	if !s.NoError(err) {
		return
	}

	stper := stepper.New(comp.Signals, comp.Steps, stepper.WithTimeout(10*time.Millisecond), stepper.WithTickDuration(2*time.Millisecond))
	stepsRun := 0
	doneCtx := RunAsync(c, comp, stper, WithPostRunFunc(func(c *engine.Context, s *step.Step) error {
		stepsRun++
		return nil
	}))

	<-doneCtx.Done()

	s.EqualError(context.Cause(doneCtx), context.DeadlineExceeded.Error())
	s.Equal(1, stepsRun)
}

//go:embed testdata/*
var testData embed.FS

func newHTMLServer() *httptest.Server {
	ser := httptest.NewServer(http.FileServerFS(testData))
	return ser
}

func newEachSelector(ls string, iter string) *v1beta1.EachSelector {
	return &v1beta1.EachSelector{
		ListSelector: ls,
		Iterator:     iter,
	}
}

func newFromDropdownComponent(each *v1beta1.EachSelector, mapper *v1beta1.Component_Dropdown_Option) *v1beta1.Component {
	return &v1beta1.Component{
		Dropdown: &v1beta1.Component_Dropdown{
			From: &v1beta1.Component_Dropdown_FromSpec{
				Selector: each,
				Mapper:   mapper,
			},
		},
	}
}

func newPromptAction(prmpt *v1beta1.Prompt) *v1beta1.Action {
	return &v1beta1.Action{
		Ask: &v1beta1.Action_PromptUser{
			Prompt: prmpt,
		},
	}
}

func newInput(selector, value string) *v1beta1.Action {
	return &v1beta1.Action{
		Input: &v1beta1.Action_Input{
			Selector: selector,
			Value:    value,
		},
	}
}

func newMouseClick(selector string) *v1beta1.Action {
	return &v1beta1.Action{
		MouseClick: &v1beta1.Action_MouseClick{
			Selector: selector,
		},
	}
}

func newNavAction(url string) *v1beta1.Action {
	return &v1beta1.Action{
		Nav: &v1beta1.Action_Nav{Addr: url},
	}
}

func TestName(t *testing.T) {
	b := []byte{}
	//_, err := os.Stdin.Read(b)

	_, err := bytes.NewBuffer(nil).Read(b)
	fmt.Println(err)
}

func TestScriptTestSuite(t *testing.T) {
	suite.Run(t, new(ScriptTestSuite))
}

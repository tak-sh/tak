package debug

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/headless/step/stepper"
	"github.com/tak-sh/tak/pkg/mocks/enginemocks"
	"testing"
	"time"
)

type StepperTestSuite struct {
	suite.Suite
}

func (s *StepperTestSuite) TestNextAndStep() {
	// -- Given
	//
	s1, _ := step.NewStep(&v1beta1.Step{Action: &v1beta1.Action{Nav: &v1beta1.Action_Nav{Addr: "derp.com"}}})
	success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
	st := NewStepper([]*step.ConditionalSignal{
		{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
	}, []*step.Step{
		s1,
	}).(*debugStepper)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	time.AfterFunc(3*time.Millisecond, func() {
		st.Step()
	})
	defer cancel()
	c := &engine.Context{Context: ctx, TemplateData: &engine.TemplateData{ScriptTemplateData: &v1beta1.ScriptTemplateData{}}}

	// -- When
	//
	h := st.Next(c)

	// -- Then
	//
	s.Equal(h.Node().Val, s1)
}

func (s *StepperTestSuite) TestNextDeadline() {
	// -- Given
	//
	st := NewStepper(nil, nil).(*debugStepper)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	c := &engine.Context{Context: ctx}

	// -- When
	//
	h := st.Next(c)

	// -- Then
	//
	s.EqualError(h.Err(), context.DeadlineExceeded.Error())
}

func (s *StepperTestSuite) TestReplay() {
	// -- Given
	//
	s1, _ := step.NewStep(&v1beta1.Step{Action: &v1beta1.Action{Nav: &v1beta1.Action_Nav{Addr: "derp1.com"}}})
	success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
	st := NewStepper([]*step.ConditionalSignal{
		{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
	}, []*step.Step{s1}, stepper.WithTickDuration(5*time.Millisecond), stepper.WithTimeout(50*time.Millisecond)).(*debugStepper)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	brow := new(enginemocks.Browser)
	brow.EXPECT().RefreshPage(mock.Anything, mock.Anything).Return(nil)
	brow.EXPECT().URL(mock.Anything).Return("", nil)
	c := &engine.Context{Context: ctx, TemplateData: &engine.TemplateData{
		ScriptTemplateData: &v1beta1.ScriptTemplateData{Browser: &v1beta1.BrowserTemplateData{}},
	}, Browser: brow}
	st.Step()
	st.Next(c)

	// -- When
	//
	st.Replay()
	actual := st.Next(c)

	// -- Then
	//
	s.Equal(actual.Node().Val, s1)
}

func (s *StepperTestSuite) TestPreviousStep() {
	// -- Given
	//
	s1, _ := step.NewStep(&v1beta1.Step{Action: &v1beta1.Action{Nav: &v1beta1.Action_Nav{Addr: "derp1.com"}}})
	s2, _ := step.NewStep(&v1beta1.Step{Action: &v1beta1.Action{Nav: &v1beta1.Action_Nav{Addr: "derp2.com"}}})
	s3, _ := step.NewStep(&v1beta1.Step{Action: &v1beta1.Action{Nav: &v1beta1.Action_Nav{Addr: "derp3.com"}}})
	success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
	st := NewStepper([]*step.ConditionalSignal{
		{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
	}, []*step.Step{s1, s2, s3}, stepper.WithTickDuration(5*time.Millisecond), stepper.WithTimeout(50*time.Millisecond)).(*debugStepper)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	brow := new(enginemocks.Browser)
	brow.EXPECT().RefreshPage(mock.Anything, mock.Anything).Return(nil)
	brow.EXPECT().URL(mock.Anything).Return("", nil)
	c := &engine.Context{Context: ctx, TemplateData: &engine.TemplateData{
		ScriptTemplateData: &v1beta1.ScriptTemplateData{Browser: &v1beta1.BrowserTemplateData{}},
	}, Browser: brow}
	st.Step()
	st.Next(c)
	st.Step()
	st.Next(c)
	st.Step()
	st.Next(c)

	// -- When
	//
	st.PreviousStep()
	actual := st.Next(c)

	// -- Then
	//
	s.Equal(actual.Node().Val, s2)
}

func (s *StepperTestSuite) TestPreviousStepAtRoot() {
	// -- Given
	//
	s1, _ := step.NewStep(&v1beta1.Step{Action: &v1beta1.Action{Nav: &v1beta1.Action_Nav{Addr: "derp1.com"}}})
	success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
	st := NewStepper([]*step.ConditionalSignal{
		{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
	}, []*step.Step{s1}, stepper.WithTickDuration(5*time.Millisecond), stepper.WithTimeout(50*time.Millisecond)).(*debugStepper)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	brow := new(enginemocks.Browser)
	brow.EXPECT().RefreshPage(mock.Anything, mock.Anything).Return(nil)
	brow.EXPECT().URL(mock.Anything).Return("", nil)
	c := &engine.Context{Context: ctx, TemplateData: &engine.TemplateData{
		ScriptTemplateData: &v1beta1.ScriptTemplateData{Browser: &v1beta1.BrowserTemplateData{}},
	}, Browser: brow}

	// -- When
	//
	st.PreviousStep()
	st.Step()
	actual := st.Next(c)

	// -- Then
	//
	s.Equal(actual.Node().Val, s1)
}

func TestStepperTestSuite(t *testing.T) {
	suite.Run(t, new(StepperTestSuite))
}

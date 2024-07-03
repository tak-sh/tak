package stepper

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/mocks/enginemocks"
	"github.com/tak-sh/tak/pkg/mocks/stepmocks"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"github.com/tak-sh/tak/pkg/utils/testutils"
	"testing"
	"time"
)

type StepperTestSuite struct {
	suite.Suite
}

func (s *StepperTestSuite) TestConstructor() {
	type test struct {
		GivenSteps []*step.Step
		Err        string
		Expected   string
	}

	tests := map[string]test{
		"simple steps": {
			GivenSteps: func() []*step.Step {
				act := new(stepmocks.Action)
				act.EXPECT().String().Return("1")

				act1 := new(stepmocks.Action)
				act1.EXPECT().String().Return("2")

				act2 := new(stepmocks.Action)
				act2.EXPECT().String().Return("3")

				return []*step.Step{
					{CompiledAction: act},
					{CompiledAction: act1},
					{CompiledAction: act2},
				}
			}(),
			Expected: "1 -> 2 -> 3",
		},
		"single branch": {
			GivenSteps: func() []*step.Step {
				act := new(stepmocks.Action)
				act.EXPECT().String().Return("1")

				act1 := testutils.NewBranchAction()
				act1.Action.EXPECT().String().Return("2")

				act2 := testutils.NewBranchAction()
				act2.Action.EXPECT().String().Return("3")

				act3 := new(stepmocks.Action)
				act3.EXPECT().String().Return("4")

				act4 := testutils.NewBranchAction()
				act4.Action.EXPECT().String().Return("5")

				return []*step.Step{
					{CompiledAction: act},
					{CompiledAction: act1},
					{CompiledAction: act2},
					{CompiledAction: act3},
					{CompiledAction: act4},
				}
			}(),
			Expected: "1 -> 2, 3, 4 -> 5",
		},
		"should error on empty": {
			GivenSteps: []*step.Step{},
			Err:        "invalid: at least 1 step required",
		},
	}

	for desc, t := range tests {
		sc := New([]*step.ConditionalSignal{
			{ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
		}, t.GivenSteps)

		s.Equal(t.Expected, sc.String(), desc)
	}
}

func (s *StepperTestSuite) TestNext() {
	type test struct {
		Stepper        Stepper
		Ctx            *engine.Context
		Expected       string
		ExpectedErr    string
		ExpectedSignal *step.ConditionalSignal
		AfterFunc      func()
	}

	tests := map[string]test{
		"runs steps to success": func() test {
			act1 := new(stepmocks.Action)

			c := &engine.Context{
				Context: context.Background(),
				TemplateData: &engine.TemplateData{
					ScriptTemplateData: &v1beta1.ScriptTemplateData{
						Browser: &v1beta1.BrowserTemplateData{},
					},
				},
			}

			act1.EXPECT().String().Return("success")
			success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
			st := New(
				[]*step.ConditionalSignal{
					{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
				},
				[]*step.Step{
					{CompiledAction: act1},
				},
			)

			return test{
				Stepper:  st,
				Ctx:      c,
				Expected: "success",
			}
		}(),
		"signals when condition met when step is ready": func() test {
			act1 := testutils.NewAction()

			c := &engine.Context{
				Context: context.Background(),
				TemplateData: &engine.TemplateData{
					ScriptTemplateData: &v1beta1.ScriptTemplateData{
						Browser: &v1beta1.BrowserTemplateData{
							Url: "derp.com",
						},
					},
				},
			}

			act1.PathNode.EXPECT().IsReady(mock.Anything).Return(true)

			success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
			expectedSig := &step.ConditionalSignal{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}}
			st := New(
				[]*step.ConditionalSignal{
					expectedSig,
				},
				[]*step.Step{
					{CompiledAction: act1},
				},
			)

			return test{
				Stepper:        st,
				Ctx:            c,
				ExpectedSignal: expectedSig,
				Expected:       "success",
			}
		}(),
		"deadlines if success signal not met": func() test {
			act1 := testutils.NewAction()

			brow := new(enginemocks.Browser)
			c := &engine.Context{
				Context: context.Background(),
				TemplateData: &engine.TemplateData{
					ScriptTemplateData: &v1beta1.ScriptTemplateData{
						Browser: &v1beta1.BrowserTemplateData{},
					},
				},
				Browser: brow,
			}

			brow.EXPECT().RefreshPage(mock.Anything, &c.TemplateData.Browser.Content).Return(nil)
			brow.EXPECT().URL(mock.Anything).Return("derp1.com", nil)

			act1.PathNode.EXPECT().IsReady(mock.Anything).Return(false)
			success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
			st := New(
				[]*step.ConditionalSignal{
					{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
				},
				[]*step.Step{
					{CompiledAction: act1},
				},
				WithTickDuration(1*time.Millisecond),
				WithTimeout(10*time.Millisecond),
			)

			return test{
				Stepper:     st,
				Ctx:         c,
				ExpectedErr: "context deadline exceeded",
			}
		}(),
		"chooses correct branch": func() test {
			act1 := testutils.NewBranchAction()
			act2 := testutils.NewBranchAction()

			brow := new(enginemocks.Browser)
			c := &engine.Context{
				Context: context.Background(),
				TemplateData: &engine.TemplateData{
					ScriptTemplateData: &v1beta1.ScriptTemplateData{
						Browser: &v1beta1.BrowserTemplateData{},
					},
				},
				Browser: brow,
			}

			brow.EXPECT().RefreshPage(mock.Anything, &c.TemplateData.Browser.Content).Return(nil)
			brow.EXPECT().URL(mock.Anything).Return("derp1.com", nil)

			act1.PathNode.EXPECT().IsReady(mock.Anything).Return(false)

			calls := 0
			act2.PathNode.EXPECT().IsReady(mock.Anything).RunAndReturn(func(data *engine.Context) bool {
				calls++
				return calls > 2
			})
			act2.Action.EXPECT().String().Return("success")

			success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
			st := New(
				[]*step.ConditionalSignal{
					{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
				},
				[]*step.Step{
					{CompiledAction: act1},
					{CompiledAction: act2},
				},
				WithTickDuration(1*time.Millisecond),
				WithTimeout(10*time.Millisecond),
			)

			return test{
				Stepper:  st,
				Ctx:      c,
				Expected: "success",
				AfterFunc: func() {
					act1.PathNode.AssertExpectations(s.T())
					act2.PathNode.AssertExpectations(s.T())
					act2.Action.AssertExpectations(s.T())
					brow.AssertExpectations(s.T())
				},
			}
		}(),
		"properly handles error signal": func() test {
			act1 := testutils.NewAction()

			brow := new(enginemocks.Browser)
			c := &engine.Context{
				Context: context.Background(),
				TemplateData: &engine.TemplateData{
					ScriptTemplateData: &v1beta1.ScriptTemplateData{
						Browser: &v1beta1.BrowserTemplateData{},
					},
				},
				Browser: brow,
			}

			brow.EXPECT().RefreshPage(mock.Anything, &c.TemplateData.Browser.Content).Return(nil)
			brow.EXPECT().URL(mock.Anything).Return("derp1.com", nil)

			act1.PathNode.EXPECT().IsReady(mock.Anything).Return(false)

			success, _ := engine.CompileTemplate("{{browser.url == 'derp.com'}}")
			errSignal, _ := engine.CompileTemplate("{{browser.url == 'derp1.com'}}")
			expectedSig := &step.ConditionalSignal{Conditional: errSignal, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_error, Message: ptr.Ptr("derp")}}
			st := New(
				[]*step.ConditionalSignal{
					{Conditional: success, ConditionalSignal: &v1beta1.ConditionalSignal{Signal: v1beta1.ConditionalSignal_success}},
					expectedSig,
				},
				[]*step.Step{
					{CompiledAction: act1},
				},
			)

			return test{
				Stepper:        st,
				Ctx:            c,
				ExpectedErr:    "derp",
				ExpectedSignal: expectedSig,
			}
		}(),
	}

	for desc, t := range tests {
		func() {
			ctx, cancel := context.WithTimeout(t.Ctx, 1*time.Second)
			t.Ctx = t.Ctx.WithContext(ctx)
			defer cancel()
			actual := t.Stepper.Next(t.Ctx)
			if actual.Err() != nil {
				s.EqualError(actual.Err(), t.ExpectedErr)
			} else {
				s.Equal(t.Expected, actual.String(), desc)
			}
			if t.ExpectedSignal != nil {
				s.Equal(t.ExpectedSignal, actual.Signal())
			}

			if t.AfterFunc != nil {
				t.AfterFunc()
			}
		}()
	}
}

func TestScannerTestSuite(t *testing.T) {
	suite.Run(t, new(StepperTestSuite))
}

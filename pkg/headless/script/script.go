package script

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/eddieowens/opts"
	"github.com/go-rod/stealth"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/action"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"github.com/tak-sh/tak/pkg/validate"
	"log/slog"
	"strconv"
	"time"
)

var desktop = device.Info{
	Name:      "Desktop",
	UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:124.0) Gecko/20100101 Firefox/124.0",
	Width:     1920,
	Height:    1080,
	Scale:     1.0,
}

type RunOpts struct {
	StartingStep int
	Store        *engine.TemplateData
	PreRun       []chromedp.Action
}

func (r RunOpts) DefaultOptions() RunOpts {
	return RunOpts{}
}

func WithPreRunActions(act ...chromedp.Action) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.PreRun = append(r.PreRun, act...)
	}
}

func WithStartingStep(i int) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.StartingStep = i
	}
}

func WithStore(st *engine.TemplateData) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.Store = st
	}
}

func Run(ctx context.Context, s *Script, str renderer.Stream, o ...opts.Opt[RunOpts]) (context.Context, error) {
	ctx, cancel := context.WithCancelCause(ctx)

	op := opts.DefaultApply(o...)
	logger := contexts.GetLogger(ctx)
	c, err := engine.NewContext(ctx, str, engine.ContextOpts{
		ScreenshotDir: s.ScreenshotsDir,
	})
	c.TemplateData = c.TemplateData.Merge(op.Store)

	chromeOpts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.UserDataDir("./data"),
	)
	execCtx, execCancel := chromedp.NewExecAllocator(ctx, chromeOpts...)

	chromeCtx, chromeCancel := chromedp.NewContext(execCtx)
	if err != nil {
		return nil, err
	}

	acts := []chromedp.Action{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(stealth.JS).Do(ctx)
			return err
		}),
		chromedp.Emulate(desktop),
	}

	acts = append(acts, op.PreRun...)
	acts = append(acts, chromedp.ActionFunc(func(ctx context.Context) error {
		c.Context = ctx
		n := len(s.Steps)

		toRedact := make([]engine.DOMDataWriter, 0, len(s.Steps))

		for i := op.StartingStep; i < n; i++ {
			v := s.Steps[i]
			logger.Info("Running action.", slog.String("action", v.Action.String()))

			sel, ok := action.AsDOMReader(v.Action)
			if ok {
				toRedact = append(toRedact, sel)
			}

			if i > 0 {
				err := chromedp.OuterHTML("html", &c.TemplateData.Browser.Content).Do(ctx)
				if err != nil {
					logger.Error("Failed to load page.", slog.String("err", err.Error()))
				}

				var u string
				err = chromedp.Location(&u).Do(ctx)
				if err != nil {
					logger.Error("Failed to get browser URL.")
				}

				if s.ScreenShotBefore {
					_, screenErr := c.Screenshot(c.Context, v.step.GetId())
					if screenErr != nil {
						logger.Error("Failed to take screenshot.", slog.String("id", v.step.GetId()), slog.String("err", screenErr.Error()))
					}
				}
			}

			if s.EventQueue != nil {
				s.EventQueue <- &ChangeStepEvent{
					Step:  v,
					Idx:   i,
					Total: n,
				}
			}

			err := runAction(c, v.Action, 10*time.Second)
			errored := err != nil
			if errored {
				logger.Error("Failed to run step.", slog.String("id", v.step.GetId()), slog.String("err", err.Error()))
			}

			if errored || s.ScreenShotAfter {
				fp, screenErr := c.Screenshot(c.Context, v.step.GetId())
				if screenErr != nil {
					logger.Error("Failed to take screenshot.", slog.String("id", v.step.GetId()), slog.String("err", screenErr.Error()))
				}
				if errored {
					return errors.Join(fmt.Errorf("failed to run step %s, see what happened here: %s", v.step.GetId(), fp), err)
				}
			}
		}
		return nil
	}))

	go func() {
		var err error
		defer func() {
			cancel(err)
			chromeCancel()
			execCancel()
		}()
		err = chromedp.Run(chromeCtx, acts...)
		return
	}()

	return ctx, nil
}

func runAction(c *engine.Context, act action.Action, to time.Duration) error {
	if _, ok := act.(*action.PromptAction); ok {
		return act.Act(c)
	}
	var toCancel context.CancelFunc
	oldCtx := c.Context
	c.Context, toCancel = context.WithTimeout(c.Context, to)
	defer func() {
		toCancel()
		c.Context = oldCtx
	}()
	err := act.Act(c)
	if errors.Is(err, context.DeadlineExceeded) {
		return except.NewTimeout("took too long")
	}
	return err
}

func New(s *v1beta1.Script) (*Script, error) {
	if len(s.Steps) == 0 {
		return nil, except.NewInvalid("at least 1 step required")
	}

	out := &Script{
		Steps:  make([]*Step, 0, len(s.Steps)),
		script: s,
	}

	for i := range s.GetSteps() {
		v := s.Steps[i]

		if v.Id == nil {
			v.Id = ptr.Ptr(strconv.Itoa(i))
		}

		st := NewStep(v)

		out.Steps = append(out.Steps, st)
	}

	err := out.Validate()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func NewStep(s *v1beta1.Step) *Step {
	return &Step{
		Action: action.New(fmt.Sprintf("step.%s", s.GetId()), s.GetAction()),
		step:   s,
	}
}

var _ grpcutils.ProtoWrapper[*v1beta1.Step] = &Step{}
var _ validate.Validator = &Step{}

type Step struct {
	Action action.Action
	step   *v1beta1.Step
}

func (s *Step) Validate() error {
	if s.step.GetAction().GetAsk() != nil && s.step.Id == nil {
		return except.NewInvalid("any step with a prompt must have an ID")
	}

	err := s.Action.Validate()
	if err != nil {
		if s.step.Id != nil {
			return errors.Join(fmt.Errorf("id %s", s.step.GetId()), err)
		}
		return err
	}

	return nil
}

func (s *Step) ToProto() *v1beta1.Step {
	return s.step
}

var _ grpcutils.ProtoWrapper[*v1beta1.Script] = &Script{}
var _ validate.Validator = &Script{}

type Script struct {
	Steps            []*Step
	ScreenShotBefore bool
	ScreenShotAfter  bool
	ScreenshotsDir   string
	EventQueue       EventQueue

	script *v1beta1.Script
}

func (c *Script) Validate() error {
	ids := map[string]int{}
	for i, v := range c.Steps {
		if v.step.Id != nil {
			id := *v.step.Id
			num, ok := ids[id]
			if ok {
				return except.NewInvalid("step #%d and #%d have the same id", num, i)
			}
			ids[id] = i
		}
	}

	for i, v := range c.Steps {
		err := v.Validate()
		if err != nil {
			return errors.Join(fmt.Errorf("step #%d", i), err)
		}
	}

	return nil
}

func (c *Script) ToProto() *v1beta1.Script {
	return c.script
}

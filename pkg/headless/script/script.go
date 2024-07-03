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
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/headless/step/stepper"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"github.com/tak-sh/tak/pkg/validate"
	"log/slog"
	"slices"
	"strconv"
)

var desktop = device.Info{
	Name:      "Desktop",
	UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:124.0) Gecko/20100101 Firefox/124.0",
	Width:     1920,
	Height:    1080,
	Scale:     1.0,
}

type RunOpts struct {
	Store       *engine.TemplateData
	PreRun      []chromedp.Action
	PostRunFunc PostRunFunc
	ChromeOpts  []chromedp.ExecAllocatorOption
}

type PostRunFunc func(c *engine.Context, s *step.Step) error

func (r RunOpts) DefaultOptions() RunOpts {
	return RunOpts{}
}

func WithChromeOpts(o ...chromedp.ExecAllocatorOption) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.ChromeOpts = append(r.ChromeOpts, o...)
	}
}

func WithPostRunFunc(f PostRunFunc) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.PostRunFunc = f
	}
}

func Run(c *engine.Context, s *Script, st stepper.Stepper, o ...opts.Opt[RunOpts]) (context.Context, error) {
	ctx, cancel := context.WithCancelCause(c.Context)

	op := opts.DefaultApply(o...)

	logger := contexts.GetLogger(ctx)
	c.TemplateData = c.TemplateData.Merge(op.Store)

	chromeOpts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		op.ChromeOpts...,
	)
	execCtx, _ := chromedp.NewExecAllocator(ctx, chromeOpts...)

	chromeCtx, _ := chromedp.NewContext(execCtx)

	acts := []chromedp.Action{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(stealth.JS).Do(ctx)
			return err
		}),
		chromedp.Emulate(desktop),
	}

	acts = append(acts, op.PreRun...)
	acts = append(acts, chromedp.ActionFunc(func(ctx context.Context) error {
		c = c.WithContext(ctx)

		for handle := st.Next(c); handle != nil; handle = st.Next(c) {
			if handle.Err() != nil {
				logger.Error("Failed to get next step.", slog.String("err", handle.Err().Error()))
				return handle.Err()
			}

			if handle.Signal() != nil {
				logger.Info("Reached success condition.", slog.Any("signal", handle.Signal().String()))
				return nil
			}

			v := handle.Node().Val
			logger.Info("Running action.", slog.String("action", v.Action.String()))

			if handle.Node().Idx > 0 {
				if s.ScreenShotBefore {
					_, screenErr := c.Screenshot(c.Context, v.GetId())
					if screenErr != nil {
						logger.Error("Failed to take screenshot.", slog.String("id", v.GetId()), slog.String("err", screenErr.Error()))
					}
				}
			}

			cont, err := s.evalStep(c, v, op)
			errored := err != nil
			if errored {
				logger.Error("Failed to run step.", slog.String("id", v.GetId()), slog.String("err", err.Error()))
			}

			if errored || s.ScreenShotAfter {
				fp, screenErr := c.Screenshot(c.Context, v.GetId())
				if screenErr != nil {
					logger.Error("Failed to take screenshot.", slog.String("id", v.GetId()), slog.String("err", screenErr.Error()))
				}
				if errored {
					return errors.Join(fmt.Errorf("failed to run step %s, see what happened here: %s", v.GetId(), fp), err)
				}
			}

			if !cont {
				logger.Info("Step signalled a completion. Stopping script.", slog.String("step", v.String()))
				return nil
			}

			err = c.RefreshPageState()
			if err != nil {
				contexts.GetLogger(ctx).Error("Failed to load page.", slog.String("err", err.Error()))
			}
		}

		return nil
	}))

	go func() {
		var err error
		defer func() {
			canErr := chromedp.Cancel(chromeCtx)
			if canErr != nil {
				logger.Error("Failed to clean up context.", slog.String("err", err.Error()))
			}
			cancel(err)
		}()
		err = chromedp.Run(chromeCtx, acts...)
		return
	}()

	return ctx, nil
}

func New(s *v1beta1.Script) (*Script, error) {
	if len(s.Steps) == 0 {
		return nil, except.NewInvalid("at least 1 step required")
	}

	idx := slices.IndexFunc(s.GetSignals(), func(signal *v1beta1.ConditionalSignal) bool {
		return signal.GetSignal() == v1beta1.ConditionalSignal_success
	})
	if idx < 0 {
		return nil, except.NewInvalid("at least 1 success condition required")
	}

	out := &Script{
		Steps:   make([]*step.Step, 0, len(s.Steps)),
		Signals: make([]*step.ConditionalSignal, len(s.Signals)),
		script:  s,
	}
	var err error
	for i, v := range s.GetSignals() {
		out.Signals[i], err = step.NewConditionalSignal(v)
		if err != nil {
			return nil, errors.Join(except.NewInvalid("script signal #%d", i), err)
		}
	}

	for i := range s.GetSteps() {
		v := s.Steps[i]

		if v.GetAction().GetBranch() != nil && v.GetId() == "" {
			return nil, errors.Join(except.NewInvalid("step #%d has a branch action but no id", i), err)
		}

		if v.Id == nil {
			v.Id = ptr.Ptr(strconv.Itoa(i))
		}

		st, err := step.NewStep(v)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("step #%d", i), err)
		}

		out.Steps = append(out.Steps, st)
	}

	err = out.Validate()
	if err != nil {
		return nil, err
	}

	return out, nil
}

var _ grpcutils.ProtoWrapper[*v1beta1.Script] = &Script{}
var _ validate.Validator = &Script{}

type Script struct {
	Steps            []*step.Step
	Signals          []*step.ConditionalSignal
	ScreenShotBefore bool
	ScreenShotAfter  bool

	script *v1beta1.Script
}

func (c *Script) Validate() error {
	ids := map[string]int{}
	for i, v := range c.Steps {
		if v.Id != nil {
			id := *v.Id
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

func (c *Script) evalStep(ctx *engine.Context, st *step.Step, op RunOpts) (cont bool, err error) {
	err = ctx.Evaluator.Eval(ctx, st)
	cont = err == nil
	if !cont {
		return false, err
	}

	if op.PostRunFunc != nil {
		err = op.PostRunFunc(ctx, st)
	}

	for _, v := range st.ConditionalSignals {
		if v.IsReady(ctx) {
			switch v.GetSignal() {
			case v1beta1.ConditionalSignal_success:
				return false, nil
			case v1beta1.ConditionalSignal_error:
				return false, except.NewAborted(v.GetMessage())
			}
		}
	}

	return
}

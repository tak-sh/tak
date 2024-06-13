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
	Store        *engine.TemplateData
	PreRun       []chromedp.Action
	StartingStep int
	Headless     bool
}

func (r RunOpts) DefaultOptions() RunOpts {
	return RunOpts{
		Headless: true,
	}
}

func WithHeadless(b bool) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.Headless = b
	}
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

func Run(c *engine.Context, s *Script, o ...opts.Opt[RunOpts]) (context.Context, error) {
	ctx, cancel := context.WithCancelCause(c.Context)

	op := opts.DefaultApply(o...)
	logger := contexts.GetLogger(ctx)
	c.TemplateData = c.TemplateData.Merge(op.Store)

	chromeOpts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserDataDir("./data"),
		chromedp.Flag("headless", op.Headless),
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
		c.Context = ctx
		n := len(s.Steps)
		for i := op.StartingStep; i < n; i++ {
			v := s.Steps[i]
			logger.Info("Running action.", slog.String("action", v.Action.String()))

			if i > 0 {
				err := c.RefreshPageState()
				if err != nil {
					logger.Error("Failed to load page.", slog.String("err", err.Error()))
				}

				err = chromedp.Location(&c.TemplateData.Browser.Url).Do(ctx)
				if err != nil {
					logger.Error("Failed to get browser URL.")
				}

				if s.ScreenShotBefore {
					_, screenErr := c.Screenshot(c.Context, v.GetId())
					if screenErr != nil {
						logger.Error("Failed to take screenshot.", slog.String("id", v.GetId()), slog.String("err", screenErr.Error()))
					}
				}
			}

			success, err := evalStep(c, s, logger, i)
			errored := err != nil
			if errored {
				logger.Error("Failed to run step.", slog.String("id", v.GetId()), slog.String("err", err.Error()))
			} else if !success {
				logger.Error("Failed to run step.")
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

	out := &Script{
		Steps:  make([]*step.Step, 0, len(s.Steps)),
		script: s,
	}

	var err error
	condSigs := make([]*step.ConditionalSignal, len(s.Signals))
	for i, v := range s.GetSignals() {
		condSigs[i], err = step.NewConditionalSignal(v)
		if err != nil {
			return nil, errors.Join(except.NewInvalid("script signal #%d", i), err)
		}
	}

	out.Decider = step.NewDecider(condSigs)

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
	ScreenShotBefore bool
	ScreenShotAfter  bool
	Decider          step.Decider

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

func gatherPaths(steps []*step.Step) []step.PathNode {
	out := make([]step.PathNode, 0, 1)
	n := len(steps)
	if n == 0 {
		return out
	}

	out = append(out, steps[0])

	// if the step contains a branching action, we should include the next non-branch step in case
	// the branch is not true.
	if _, ok := steps[0].CompiledAction.(*step.BranchAction); ok && n > 1 {
		for i := 1; i < len(steps); i++ {
			v := steps[i]
			out = append(out, v)
			if _, isBranch := v.CompiledAction.(*step.BranchAction); !isBranch {
				break
			}
		}
	}

	return out
}

func evalStep(c *engine.Context, s *Script, logger *slog.Logger, i int) (bool, error) {
	v := s.Steps[i]

	decided, stop, err := s.Decider.ChoosePath(c, time.Second*10, gatherPaths(s.Steps[i:])...)
	if err != nil {
		logger.Error("Failed to choose a path.", slog.String("err", err.Error()))
		return false, err
	}

	if stop {
		return false, nil
	}

	switch decided.(type) {
	case *step.Step:
		err = c.Evaluator.Eval(c, v)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		logger.Error("Unknown path", slog.Any("path", decided), slog.String("step", v.GetId()))
		return false, except.NewInternal("unsure on how to continue")
	}
}

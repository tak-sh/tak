package script

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/go-rod/stealth"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/headless/action"
	"github.com/tak-sh/tak/pkg/internal/grpcutils"
	"github.com/tak-sh/tak/pkg/internal/ptr"
	"github.com/tak-sh/tak/pkg/validate"
	"log/slog"
	"strconv"
)

var desktop = device.Info{
	Name:      "Desktop",
	UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:124.0) Gecko/20100101 Firefox/124.0",
	Width:     1920,
	Height:    1080,
	Scale:     1.0,
}

func Run(c *headless.Context, s *Script) error {
	logger := contexts.GetLogger(c)

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
	)
	ctx, execCancel := chromedp.NewExecAllocator(c, opts...)
	defer execCancel()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	acts := []chromedp.Action{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(stealth.JS).Do(ctx)
			return err
		}),
		chromedp.Emulate(desktop),
		chromedp.ActionFunc(func(ctx context.Context) error {
			c.Context = ctx

			var pageContent string
			for i := range s.Steps {
				v := s.Steps[i]
				logger.Info("Running action.", slog.String("action", v.Action.String()))

				if i > 0 {
					err := chromedp.OuterHTML("html", &pageContent).Do(ctx)
					if err != nil {
						logger.Error("Failed to load page.", slog.String("err", err.Error()))
					} else {
						c.Store.Set("page", pageContent)
					}
				}

				err := v.Action.Act(c)
				errored := err != nil
				if errored {
					logger.Error("Failed to run step.", slog.String("id", v.step.GetId()), slog.String("err", err.Error()))
				}

				if errored || s.Debug {
					screenErr := c.Screenshot(c.Context, v.step.GetId())
					if screenErr != nil {
						logger.Error("Failed to take screenshot.", slog.String("id", v.step.GetId()), slog.String("err", screenErr.Error()))
					}
					if errored {
						return errors.Join(fmt.Errorf("failed to run step %s", v.step.GetId()), err)
					}
				}
			}
			return nil
		}),
	}

	err := chromedp.Run(ctx, acts...)
	if err != nil {
		return err
	}

	return nil
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
		Action: action.New(s.GetId(), s.GetAction()),
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
	Steps []*Step
	Debug bool

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

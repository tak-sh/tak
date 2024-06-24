package ui

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eddieowens/opts"
	"github.com/tak-sh/tak/generated/go/api/account/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/utils/stringutils"
	"io"
	"log/slog"
)

type BubbleUIOpts struct {
	OnPromptFunc OnPromptFunc
}

func WithOnPrompt(f OnPromptFunc) opts.Opt[BubbleUIOpts] {
	return func(b *BubbleUIOpts) {
		b.OnPromptFunc = f
	}
}

type OnPromptFunc func(id string)

func (b BubbleUIOpts) DefaultOptions() BubbleUIOpts {
	return BubbleUIOpts{}
}

func NewBubbleUI(account *v1beta1.Account, str renderer.Stream, eq engine.EventQueue, o ...opts.Opt[BubbleUIOpts]) UI {
	op := opts.DefaultApply(o...)

	return &bubble{
		Account:      account,
		Stream:       str,
		ScriptEvents: eq,
		Op:           op,
	}
}

type bubble struct {
	Account      *v1beta1.Account
	Stream       renderer.Stream
	ScriptEvents engine.EventQueue
	Op           BubbleUIOpts
}

func (u *bubble) Start(ctx context.Context, r io.Reader, w io.Writer) (context.Context, error) {
	app := NewApp(func(s *SubmitEvent) {
		u.Stream.Respond(&renderer.Response{
			ID:    s.ID,
			Value: s.Val,
		})
	}, fmt.Sprintf("Adding your %s account...", u.Account.GetMetadata().GetName()))

	p := tea.NewProgram(
		app,
		tea.WithContext(ctx),
		tea.WithAltScreen(),
		tea.WithInput(r),
		tea.WithOutput(w),
	)
	logger := contexts.GetLogger(ctx)

	ctx, cancel := context.WithCancelCause(ctx)

	go func() {
		var err error
		defer func() {
			cancel(err)
		}()
		_, err = p.Run()
		if err != nil {
			logger.Error("Failed to run program.", slog.String("err", err.Error()))
		}
	}()
	go func() {
		var err error
		defer func() {
			cancel(err)
		}()

		for {
			select {
			case <-ctx.Done():
				logger.Info("Received done on context chan.")
				return
			case r, ok := <-u.Stream.RenderQueue():
				if !ok {
					logger.Info("Render queue closed.")
					return
				}
				p.Send(SetChildrenMsg{Models: []tea.Model{r}})
				if u.Op.OnPromptFunc != nil {
					u.Op.OnPromptFunc(r.GetId())
				}
			case e, ok := <-u.ScriptEvents:
				if !ok {
					logger.Info("Event queue has closed.")
					return
				}

				switch t := e.(type) {
				case *engine.NextInstructionEvent:
					switch in := t.Instruction.(type) {
					case *step.Step:
						if _, ok := in.CompiledAction.(*step.PromptAction); !ok {
							p.Send(UpdateProgressMsg{Msg: stringutils.Capitalize(in.CompiledAction.String())})
						}
					}
				}
			}
		}
	}()

	return ctx, nil
}

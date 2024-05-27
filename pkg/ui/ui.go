package ui

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tak-sh/tak/generated/go/api/account/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/headless/action"
	"github.com/tak-sh/tak/pkg/headless/script"
	"github.com/tak-sh/tak/pkg/internal/stringutils"
	"github.com/tak-sh/tak/pkg/renderer"
	"io"
	"log/slog"
)

type UI interface {
	Start(ctx context.Context, r io.Reader, w io.Writer) (context.Context, error)
}

func NewBubbleUI(account *v1beta1.Account, str renderer.Stream, eq script.EventQueue) UI {
	return &ui{
		Account:      account,
		Stream:       str,
		ScriptEvents: eq,
	}
}

type ui struct {
	Account      *v1beta1.Account
	Stream       renderer.Stream
	ScriptEvents script.EventQueue
}

func (u *ui) Start(ctx context.Context, r io.Reader, w io.Writer) (context.Context, error) {
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
			case e, ok := <-u.ScriptEvents:
				if !ok {
					logger.Info("Event queue has closed.")
					return
				}

				switch t := e.(type) {
				case *script.ChangeStepEvent:
					if _, ok := t.Step.Action.(*action.PromptAction); !ok {
						p.Send(UpdateProgressMsg{Msg: stringutils.Capitalize(t.Step.Action.String())})
					}
				}
			}
		}
	}()

	return ctx, nil
}

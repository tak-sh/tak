package ui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tak-sh/tak/pkg/contexts"
	"log/slog"
	"strings"
)

func Run(ctx context.Context, p *tea.Program) (context.Context, error) {
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

	return ctx, nil
}

func Indent(str string, n int) string {
	ss := strings.Split(str, "\n")
	for i, v := range ss {
		str := make([]string, 0, n+1)
		for i := 0; i < n; i++ {
			str = append(str, "")
		}
		str = append(str, v)
		ss[i] = strings.Join(str, " ")
	}
	return strings.Join(ss, "\n")
}

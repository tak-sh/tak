package contexts

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	v, _ := ctx.Value(loggerKey{}).(*slog.Logger)
	if v == nil {
		v = slog.Default()
	}
	return v
}

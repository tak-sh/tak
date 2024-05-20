package ui

import "context"

type UI interface {
	Start(ctx context.Context) (context.Context, error)
}

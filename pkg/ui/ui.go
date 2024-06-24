package ui

import (
	"context"
	"io"
)

type UI interface {
	Start(ctx context.Context, r io.Reader, w io.Writer) (context.Context, error)
}

package ui

import (
	"bytes"
	"context"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/internal/ptr"
	"github.com/tak-sh/tak/pkg/prompt"
	"io"
	"log/slog"
	"strings"
)

func NewReadWriterUI(r io.Reader, w io.Writer, stream prompt.Stream) *ReadWriter {
	out := &ReadWriter{
		R:      r,
		W:      w,
		Stream: stream,
	}

	return out
}

var _ UI = &ReadWriter{}

type ReadWriter struct {
	R      io.Reader
	W      io.Writer
	Stream prompt.Stream
}

func (s *ReadWriter) Start(ctx context.Context) (context.Context, error) {
	ctx, cancel := context.WithCancelCause(ctx)

	go func() {
		var err error
		logger := contexts.GetLogger(ctx)
		defer func() {
			logger.Info("Shutting down simple UI.", slog.String("err", err.Error()))
			cancel(err)
		}()
		for {
			readBuffer := make([]byte, 200)
			select {
			case <-ctx.Done():
				err = context.Cause(ctx)
				return
			case p, ok := <-s.Stream.Prompts():
				if !ok {
					return
				}
				b := strings.Builder{}
				b.WriteString(p.GetTitle())
				b.WriteString(":\n")
				if p.Description != nil {
					b.WriteString(*p.Description)
					b.WriteString("\n")
				}

				if dd := p.Component.GetDropdown(); dd != nil {
					b.WriteString("Options:\n")
					for _, v := range dd.GetOptions() {
						b.WriteString(v.Value)
						b.WriteString("\n")
					}
				}

				_, err = s.W.Write([]byte(b.String()))
				if err != nil {
					logger.Error("Failed to write prompt to user")
					return
				}

				n, err := s.R.Read(readBuffer)
				if err != nil {
					logger.Error("Failed to read user input.", slog.String("err", err.Error()))
				} else {
					readBuffer = bytes.TrimSuffix(readBuffer[:n], []byte("\n"))
					s.Stream.SendResponse(&v1beta1.Value{
						Str: ptr.Ptr(string(readBuffer)),
					})
				}
			}
		}
	}()

	return ctx, nil
}

package prompt

import (
	"context"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
)

type Stream interface {
	// SendPrompt sends p to the user and awaits a response. This
	// method is blocking.
	SendPrompt(ctx context.Context, p *v1beta1.Prompt) (*v1beta1.Value, error)
	Prompts() <-chan *v1beta1.Prompt
	SendResponse(v *v1beta1.Value)
}

func NewStream() Stream {
	return &stream{
		requests:  make(chan *v1beta1.Prompt, 10),
		responses: make(chan *v1beta1.Value, 10),
	}
}

var _ Stream = &stream{}

type stream struct {
	requests  chan *v1beta1.Prompt
	responses chan *v1beta1.Value
}

func (c stream) SendResponse(v *v1beta1.Value) {
	c.responses <- v
}

func (c stream) SendPrompt(ctx context.Context, p *v1beta1.Prompt) (*v1beta1.Value, error) {
	c.requests <- p

	select {
	case <-ctx.Done():
		return nil, context.Cause(ctx)
	case v, ok := <-c.responses:
		if !ok {
			return nil, except.NewAborted("user cancelled prompt")
		}
		return v, nil
	}
}

func (c stream) Prompts() <-chan *v1beta1.Prompt {
	return c.requests
}

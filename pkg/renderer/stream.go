package renderer

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
)

type Stream interface {
	// Render renders the Model to the user. This method is blocking.
	// until the user has made a decision.
	Render(ctx context.Context, p Model) (*Response, error)
	RenderQueue() <-chan Model
	Respond(v *Response)
}

type Model interface {
	tea.Model
	GetId() string
}

type Response struct {
	ID    string
	Value *v1beta1.Value
}

func NewStream() Stream {
	return &stream{
		requests:  make(chan Model, 10),
		responses: make(chan *Response, 10),
	}
}

var _ Stream = &stream{}

type stream struct {
	requests  chan Model
	responses chan *Response
}

func (c stream) Respond(v *Response) {
	c.responses <- v
}

func (c stream) Render(ctx context.Context, p Model) (*Response, error) {
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

func (c stream) RenderQueue() <-chan Model {
	return c.requests
}

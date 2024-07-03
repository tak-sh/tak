package component

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"github.com/tak-sh/tak/pkg/validate"
)

type Props struct {
	ID          string
	Title       string
	Description string
}

type Component interface {
	validate.Validator
	grpcutils.ProtoWrapper[*v1beta1.Component]
	Render(c *engine.Context, p *Props) renderer.Model
}

func New(c *v1beta1.Component) (Component, error) {
	if i := c.GetInput(); i != nil {
		return NewInput(i)
	} else if i := c.GetDropdown(); i != nil {
		return NewDropdown(i)
	}

	return nil, except.NewInvalid("blank components not allowed")
}

type OnSubmitMsg struct {
	Id    string
	Value *v1beta1.Value
}

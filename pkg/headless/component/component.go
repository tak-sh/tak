package component

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/internal/grpcutils"
	"github.com/tak-sh/tak/pkg/validate"
)

const (
	PageKey = "page"
)

type Rendered interface {
	Render(ctx *headless.Context) (Component, error)
}

type Component interface {
	validate.Validator
	Rendered
	grpcutils.ProtoWrapper[*v1beta1.Component]
}

func New(c *v1beta1.Component) Component {
	if i := c.GetInput(); i != nil {
		return NewInput(i)
	} else if i := c.GetDropdown(); i != nil {
		return NewDropdown(i)
	} else {
		return &NoOp{}
	}
}

package component

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/internal/grpcutils"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/validate"
)

const (
	PageKey = "page"
)

type Props struct {
	ID          string
	Title       string
	Description string
}

type Component interface {
	validate.Validator
	grpcutils.ProtoWrapper[*v1beta1.Component]
	Render(c *headless.Context, p *Props) renderer.Model
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

// SyncStateMsg is sent when the program has detected the user has finished populating a field.
type SyncStateMsg func(id string, v *v1beta1.Value)

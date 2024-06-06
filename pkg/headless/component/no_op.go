package component

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/renderer"
)

var _ Component = &NoOp{}

type NoOp struct {
}

func (n *NoOp) Render(_ *engine.Context, _ *Props) renderer.Model {
	return nil
}

func (n *NoOp) ToProto() *v1beta1.Component {
	return &v1beta1.Component{}
}

func (n *NoOp) Validate() error {
	return nil
}

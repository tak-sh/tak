package component

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless"
)

var _ Component = &NoOp{}

type NoOp struct {
}

func (n *NoOp) ToProto() *v1beta1.Component {
	return &v1beta1.Component{}
}

func (n *NoOp) Render(_ *headless.Context) (Component, error) {
	return n, nil
}

func (n *NoOp) Validate() error {
	return nil
}

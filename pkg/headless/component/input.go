package component

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless"
)

func NewInput(i *v1beta1.Component_Input) *Input {
	out := &Input{
		comp: i,
	}

	return out
}

var _ Component = &Input{}

type Input struct {
	comp *v1beta1.Component_Input
}

func (i *Input) ToProto() *v1beta1.Component {
	return &v1beta1.Component{Input: i.comp}
}

func (i *Input) Render(_ *headless.Context) (Component, error) {
	return i, nil
}

func (i *Input) Validate() error {
	return nil
}

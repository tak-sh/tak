package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
)

func NewInput(id string, d *v1beta1.Action_Input) (*Input, error) {
	out := &Input{
		Action_Input: d,
		ID:           id,
	}

	var err error
	out.ValueTemp, err = engine.CompileTemplate(d.Value)
	if err != nil {
		return nil, errors.Join(except.NewInvalid("invalid template value %s", d.Value), err)
	}

	return out, nil
}

var _ Action = &Input{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_Input] = &Input{}
var _ engine.PathNode = &Input{}

type Input struct {
	ID string
	*v1beta1.Action_Input
	ValueTemp *engine.TemplateRenderer
}

func (i *Input) GetId() string {
	return i.ID
}

func (i *Input) IsReady(c *engine.Context) bool {
	return c.Browser.Exists(c.Context, i.GetSelector())
}

func (i *Input) Validate() error {
	if i.GetSelector() == "" {
		return except.NewInvalid("a selector is required")
	}
	return nil
}

func (i *Input) String() string {
	return fmt.Sprintf("inputting %s into %s", i.GetValue(), i.GetSelector())
}

func (i *Input) Act(c *engine.Context) error {
	val := i.ValueTemp.Render(c.TemplateData)
	return c.Browser.WriteInput(c, i.GetSelector(), val)
}

func (i *Input) GetID() string {
	return i.ID
}

func (i *Input) ToProto() *v1beta1.Action_Input {
	return i.Action_Input
}

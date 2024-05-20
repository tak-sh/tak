package action

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/internal/grpcutils"
)

func NewInput(id string, d *v1beta1.Action_Input) *Input {
	out := &Input{
		Action_Input: d,
		ID:           id,
	}

	return out
}

var _ Action = &Input{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_Input] = &Input{}

type Input struct {
	ID string
	*v1beta1.Action_Input
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

func (i *Input) Act(c *headless.Context) error {
	val := c.RenderTemplate(i.GetValue())
	return chromedp.SendKeys(i.GetSelector(), val).Do(c)
}

func (i *Input) GetID() string {
	return i.ID
}

func (i *Input) ToProto() *v1beta1.Action_Input {
	return i.Action_Input
}

package action

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
)

func NewInput(id string, d *v1beta1.Action_Input) *Input {
	out := &Input{
		Action_Input: d,
		ID:           id,
		Query:        engine.StringSelector(d.Selector),
	}

	return out
}

var _ Action = &Input{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_Input] = &Input{}
var _ engine.DOMDataWriter = &Input{}

type Input struct {
	ID string
	*v1beta1.Action_Input
	Query engine.DOMQuery
}

func (i *Input) GetQueries() []engine.DOMQuery {
	return []engine.DOMQuery{i.Query}
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
	val := c.TemplateData.Render(i.GetValue())
	return chromedp.SendKeys(i.GetSelector(), val).Do(c)
}

func (i *Input) GetID() string {
	return i.ID
}

func (i *Input) ToProto() *v1beta1.Action_Input {
	return i.Action_Input
}

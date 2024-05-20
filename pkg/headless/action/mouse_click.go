package action

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/internal/grpcutils"
)

func NewMouseClick(id string, d *v1beta1.Action_MouseClick) *MouseClick {
	out := &MouseClick{
		Action_MouseClick: d,
		ID:                id,
	}

	return out
}

var _ Action = &MouseClick{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_MouseClick] = &MouseClick{}

type MouseClick struct {
	*v1beta1.Action_MouseClick
	ID string
}

func (m *MouseClick) Validate() error {
	if m.GetSelector() == "" {
		return except.NewInvalid("a selector is required")
	}
	return nil
}

func (m *MouseClick) String() string {
	click := "clicking"
	if m.GetDouble() {
		click = "double clicking"
	}
	return fmt.Sprintf("%s on %s", click, m.GetSelector())
}

func (m *MouseClick) Act(c *headless.Context) error {
	if m.GetDouble() {
		return chromedp.DoubleClick(m.GetSelector()).Do(c)
	} else {
		return chromedp.Click(m.GetSelector()).Do(c)
	}
}

func (m *MouseClick) GetID() string {
	return m.ID
}

func (m *MouseClick) ToProto() *v1beta1.Action_MouseClick {
	return m.Action_MouseClick
}

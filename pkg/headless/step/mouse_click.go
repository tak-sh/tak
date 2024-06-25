package step

import (
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
)

func NewMouseClick(id string, d *v1beta1.Action_MouseClick) (*MouseClick, error) {
	out := &MouseClick{
		Action_MouseClick: d,
		ID:                id,
		Query:             engine.StringSelector(d.Selector),
	}

	var err error
	out.SelectorTemp, err = engine.CompileTemplate(d.Selector)
	if err != nil {
		return nil, errors.Join(except.NewInvalid("invalid template selector %s", d.Selector), err)
	}

	return out, nil
}

var _ Action = &MouseClick{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_MouseClick] = &MouseClick{}
var _ engine.PathNode = &MouseClick{}

type MouseClick struct {
	*v1beta1.Action_MouseClick
	ID           string
	Query        engine.DOMQuery
	SelectorTemp *engine.TemplateRenderer
}

func (m *MouseClick) GetId() string {
	return m.ID
}

func (m *MouseClick) IsReady(c *engine.Context) bool {
	return c.Browser.Exists(c.Context, m.SelectorTemp.Render(c.TemplateData))
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

func (m *MouseClick) Act(c *engine.Context) error {
	sel := m.SelectorTemp.Render(c.TemplateData)
	if m.GetDouble() {
		return chromedp.DoubleClick(sel).Do(c)
	} else {
		return chromedp.Click(sel).Do(c)
	}
}

func (m *MouseClick) GetID() string {
	return m.ID
}

func (m *MouseClick) ToProto() *v1beta1.Action_MouseClick {
	return m.Action_MouseClick
}

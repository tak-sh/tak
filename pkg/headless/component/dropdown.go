package component

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"strings"
)

func NewDropdown(d *v1beta1.Component_Dropdown) *Dropdown {
	out := &Dropdown{
		comp: d,
	}

	return out
}

var _ Component = &Dropdown{}

type Dropdown struct {
	comp *v1beta1.Component_Dropdown
}

func (d *Dropdown) ToProto() *v1beta1.Component {
	return &v1beta1.Component{Dropdown: d.comp}
}

func (d *Dropdown) Render(ctx *headless.Context) (Component, error) {
	cl := proto.Clone(d.comp).(*v1beta1.Component_Dropdown)
	logger := contexts.GetLogger(ctx)

	for _, v := range cl.GetOptions() {
		v.Value = ctx.RenderTemplate(v.Value)
	}

	if d.comp.From != nil {
		listSelector := ctx.RenderTemplate(d.comp.From.GetListSelector())
		v := ctx.Store.Get(PageKey)
		if v != nil {
			raw := v.(string)
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
			if err != nil {
				logger.Error("Failed to render dropdown from field dur to bad HTML doc.",
					slog.String("err", err.Error()),
					slog.String("field", d.comp.From.ListSelector),
					slog.String("doc", raw),
				)
				return NewDropdown(cl), nil
			}

			sel := doc.Find(listSelector)
			sel.Each(func(i int, selection *goquery.Selection) {
				sel = selection.Find(d.comp.From.Iterator)
				if sel == nil {
					return
				}
				text := sel.Text()

				cl.Options = append(cl.Options, &v1beta1.Component_Dropdown_Option{Value: text})
			})
		}
	}

	return NewDropdown(cl), nil
}

func (d *Dropdown) Validate() error {
	if len(d.comp.GetOptions()) == 0 && d.comp.From == nil {
		return except.NewInvalid("at least one option or a from field is required")
	} else if d.comp.From != nil {
		if d.comp.From.ListSelector == "" || d.comp.From.Iterator == "" {
			return except.NewInvalid("both the list_selector and iterator fields are required")
		}
	}

	for _, v := range d.comp.GetOptions() {
		err := validateDropdownComponentOption(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateDropdownComponentOption(_ *v1beta1.Component_Dropdown_Option) error {
	return nil
}

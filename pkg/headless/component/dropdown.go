package component

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/internal/ptr"
	"github.com/tak-sh/tak/pkg/renderer"
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

func (d *Dropdown) Render(ctx *headless.Context, props *Props) renderer.Model {
	cl := proto.Clone(d.comp).(*v1beta1.Component_Dropdown)
	logger := contexts.GetLogger(ctx)

	for _, v := range cl.GetOptions() {
		v.Value = ctx.Store.Render(v.Value)
	}

	if d.comp.From != nil {
		listSelector := ctx.Store.Render(d.comp.From.GetSelector())
		v := ctx.Store.Get(PageKey)
		if v != nil {
			raw := v.(string)
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
			if err != nil {
				logger.Error("Failed to render dropdown from field dur to bad HTML doc.",
					slog.String("err", err.Error()),
					slog.String("field", d.comp.From.Selector),
					slog.String("doc", raw),
				)
				return newDropdownModel(cl, props)
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

	for _, mer := range cl.GetMerge() {
		for _, opt := range cl.GetOptions() {
			ifVal := ctx.Store.Merge(headless.Store{"option": opt}).Render(mer.GetIf())
			if headless.IsTruthy(ifVal) {
				proto.Merge(opt, mer.GetOption())
			}
		}
	}

	return newDropdownModel(cl, props)
}

func (d *Dropdown) Validate() error {
	if len(d.comp.GetOptions()) == 0 && d.comp.From == nil {
		return except.NewInvalid("at least one option or a from field is required")
	} else if d.comp.From != nil {
		if d.comp.From.Selector == "" || d.comp.From.Iterator == "" {
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

var _ renderer.Model = &DropdownModel{}

func newDropdownModel(comp *v1beta1.Component_Dropdown, props *Props) *DropdownModel {
	items := make([]list.Item, 0, len(comp.GetOptions()))
	for i := range comp.GetOptions() {
		v := comp.Options[i]
		items = append(items, &dropdownItem{comp: v})
	}

	return &DropdownModel{
		Props: props,
		List:  list.New(items, list.NewDefaultDelegate(), 0, 0),
		Comp:  comp,
	}
}

type DropdownModel struct {
	Props *Props
	List  list.Model
	Comp  *v1beta1.Component_Dropdown
}

func (d *DropdownModel) Init() tea.Cmd {
	return nil
}

func (d *DropdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	d.List, cmd = d.List.Update(msg)

	switch t := msg.(type) {
	case SyncStateMsg:
		v, ok := d.List.SelectedItem().(*dropdownItem)
		if ok {
			t(d.Props.ID, &v1beta1.Value{Str: ptr.Ptr(v.comp.Value)})
		}
	}

	return d, cmd
}

func (d *DropdownModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, d.Props.Title, d.Props.Description, d.List.View())
}

var _ list.Item = &dropdownItem{}

type dropdownItem struct {
	comp *v1beta1.Component_Dropdown_Option
}

func (d *dropdownItem) FilterValue() string {
	return d.comp.GetValue()
}

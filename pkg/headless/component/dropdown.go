package component

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"golang.org/x/net/html"
	"google.golang.org/protobuf/proto"
	"io"
	"log/slog"
	"strings"
)

func NewDropdown(d *v1beta1.Component_Dropdown) (*Dropdown, error) {
	out := &Dropdown{
		Component_Dropdown: d,
		Query:              engine.NewEachSelector(d.GetFrom().GetSelector()),
		OptionValueTemps:   make([]*engine.TemplateRenderer, 0, len(d.Options)),
		OptionMergeIfTemps: make([]*engine.TemplateRenderer, 0, len(d.Merge)),
	}

	for i, v := range d.GetOptions() {
		op, err := engine.CompileTemplate(v.Value)
		if err != nil {
			return nil, errors.Join(except.NewInvalid("%s is an invalid template for option #%d", v.Value, i), err)
		}
		out.OptionValueTemps = append(out.OptionValueTemps, op)
	}

	var err error
	if val := d.GetFrom().GetMapper().GetValue(); val != "" {
		out.MapperOptionValueTemp, err = engine.CompileTemplate(val)
		if err != nil {
			return nil, errors.Join(except.NewInvalid("%s is an invalid template for mapper value field", val), err)
		}
	}

	if val := d.GetFrom().GetMapper().GetText(); val != "" {
		out.MapperOptionTextTemp, err = engine.CompileTemplate(val)
		if err != nil {
			return nil, errors.Join(except.NewInvalid("%s is an invalid template for mapper text field", val), err)
		}
	}

	for i, mer := range d.GetMerge() {
		ifStr := "true"
		if mer.If != nil {
			ifStr = *mer.If
		}
		op, err := engine.CompileTemplate(ifStr)
		if err != nil {
			return nil, errors.Join(except.NewInvalid("%s is an invalid template for the if field in merged option #%d", ifStr, i), err)
		}
		out.OptionMergeIfTemps = append(out.OptionMergeIfTemps, op)
	}

	if d.GetFrom().GetSelector() != nil {
		out.FromListSelectorTemp, err = engine.CompileTemplate(d.GetFrom().GetSelector().GetListSelector())
		if err != nil {
			return nil, errors.Join(except.NewInvalid("%s is an invalid template for the list_selector field", d.From.Selector.ListSelector), err)
		}
	}

	return out, nil
}

var _ Component = &Dropdown{}
var _ engine.DOMDataWriter = &Dropdown{}

type Dropdown struct {
	*v1beta1.Component_Dropdown
	Query engine.DOMQuery

	OptionValueTemps      []*engine.TemplateRenderer
	MapperOptionValueTemp *engine.TemplateRenderer
	MapperOptionTextTemp  *engine.TemplateRenderer
	OptionMergeIfTemps    []*engine.TemplateRenderer
	FromListSelectorTemp  *engine.TemplateRenderer
}

func (d *Dropdown) GetQueries() []engine.DOMQuery {
	return []engine.DOMQuery{d.Query}
}

func (d *Dropdown) ToProto() *v1beta1.Component {
	return &v1beta1.Component{Dropdown: d.Component_Dropdown}
}

func (d *Dropdown) Render(ctx *engine.Context, props *Props) renderer.Model {
	cl := proto.Clone(d.Component_Dropdown).(*v1beta1.Component_Dropdown)
	logger := contexts.GetLogger(ctx)

	for i, v := range cl.GetOptions() {
		v.Value = d.OptionValueTemps[i].Render(ctx.TemplateData)
	}

	if d.GetFrom().GetSelector() != nil {
		from := proto.Clone(d.GetFrom()).(*v1beta1.Component_Dropdown_FromSpec)
		from.Selector.ListSelector = d.FromListSelectorTemp.Render(ctx.TemplateData)
		raw := ctx.TemplateData.GetBrowser().GetContent()
		if raw != "" {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
			if err != nil {
				logger.Error("Failed to render dropdown from field due to bad HTML doc.",
					slog.String("err", err.Error()),
					slog.String("field", d.From.Selector.ListSelector),
					slog.String("doc", raw),
				)
				return newDropdownModel(cl, props)
			}

			eles := engine.NewEachSelector(from.GetSelector()).Query(doc.Selection)
			for _, sel := range eles {
				st := addDocToStore(ctx.TemplateData, sel)

				opt := proto.Clone(from.Mapper).(*v1beta1.Component_Dropdown_Option)

				opt.Value = d.MapperOptionValueTemp.Render(st)
				opt.Text = ptr.PtrOrNil(d.MapperOptionTextTemp.Render(st))

				cl.Options = append(cl.Options, opt)
			}
		}
	}

	for i, mer := range cl.GetMerge() {
		for _, opt := range cl.GetOptions() {
			tempData := ctx.TemplateData.Merge(&engine.TemplateData{
				ScriptTemplateData: &v1beta1.ScriptTemplateData{
					Option: opt,
				},
			})
			if engine.IsTruthy(d.OptionMergeIfTemps[i].Render(tempData)) {
				proto.Merge(opt, mer.GetOption())
			}
		}
	}

	return newDropdownModel(cl, props)
}

func (d *Dropdown) Validate() error {
	if len(d.GetOptions()) == 0 && d.From == nil {
		return except.NewInvalid("at least one option or a from field is required")
	} else if d.From != nil {
		if d.From.Selector.ListSelector == "" || d.From.Selector.Iterator == "" || d.From.Selector == nil {
			return except.NewInvalid("both the extract and mapper fields are required")
		}
	}

	for _, v := range d.GetOptions() {
		err := validateDropdownComponentOption(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func addDocToStore(st *engine.TemplateData, sel *goquery.Selection) *engine.TemplateData {
	if sel == nil {
		return st
	}

	if len(sel.Nodes) == 0 {
		return st
	}

	return st.Merge(&engine.TemplateData{
		ScriptTemplateData: &v1beta1.ScriptTemplateData{
			Element: nodeToTemplate(sel.Nodes[0]),
		},
	})
}

func nodeToTemplate(node *html.Node) *v1beta1.HTMLNodeTemplateData {
	out := &v1beta1.HTMLNodeTemplateData{
		Attrs:   make(map[string]*v1beta1.HTMLNodeTemplateData_Attribute),
		Element: node.DataAtom.String(),
	}

	if node.FirstChild != nil {
		out.Data = node.FirstChild.Data
	}

	for _, v := range node.Attr {
		out.Attrs[v.Key] = &v1beta1.HTMLNodeTemplateData_Attribute{
			Val:       v.Val,
			Namespace: v.Namespace,
		}
	}

	return out
}

func validateDropdownComponentOption(_ *v1beta1.Component_Dropdown_Option) error {
	return nil
}

var _ renderer.Model = &DropdownModel{}

func newDropdownModel(comp *v1beta1.Component_Dropdown, props *Props) *DropdownModel {
	items := make([]list.Item, 0, len(comp.GetOptions()))
	displayIdx := 1
	for i := range comp.GetOptions() {
		v := comp.Options[i]
		if v.GetHidden() {
			continue
		}

		item := &dropdownItem{comp: v, idx: i}
		if !v.GetDisabled() {
			item.displayIdx = displayIdx
			displayIdx++
		}

		items = append(items, item)
	}

	li := list.New(items, &dropdownItemDelegate{}, 0, len(items)*2)

	return &DropdownModel{
		Props: props,
		List:  li,
		Comp:  comp,
	}
}

type DropdownModel struct {
	Props *Props
	List  list.Model
	Comp  *v1beta1.Component_Dropdown
}

func (d *DropdownModel) GetId() string {
	return d.Props.ID
}

func (d *DropdownModel) Init() tea.Cmd {
	if len(d.Comp.Options) > 0 && d.Comp.Options[d.List.Index()].GetDisabled() {
		d.List.CursorDown()
	}
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
	case tea.KeyMsg:
		if key.Matches(t, list.DefaultKeyMap().CursorUp, list.DefaultKeyMap().CursorDown) {
			idx := d.List.Index()
			if !d.Comp.Options[idx].GetDisabled() {
				break
			}
			halt := idx == len(d.Comp.Options)-1 || idx == 0
			up := key.Matches(t, list.DefaultKeyMap().CursorUp)
			var direc func()
			if halt {
				if up {
					direc = d.List.CursorDown
				} else {
					direc = d.List.CursorUp
				}
			} else if up {
				direc = d.List.CursorUp
			} else {
				direc = d.List.CursorDown
			}
			direc()
		}
	}

	return d, cmd
}

func (d *DropdownModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, d.Props.Title, d.Props.Description, d.List.View())
}

var _ list.ItemDelegate = &dropdownItemDelegate{}

type dropdownItemDelegate struct {
}

func (d *dropdownItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(*dropdownItem)
	if !ok {
		return
	}

	fn := DropdownItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedDropdownItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(i.getText()))
}

func (d *dropdownItemDelegate) Height() int {
	return 1
}

func (d *dropdownItemDelegate) Spacing() int {
	return 0
}

func (d *dropdownItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

var _ list.Item = &dropdownItem{}

type dropdownItem struct {
	idx        int
	displayIdx int
	comp       *v1beta1.Component_Dropdown_Option
}

func (d *dropdownItem) getText() string {
	var str string
	if d.comp.GetDisabled() {
		if d.comp.Text != nil {
			str = *d.comp.Text
		} else {
			str = d.comp.Value
		}
	} else {
		text := d.comp.Value
		if d.comp.Text != nil {
			text = *d.comp.Text
		}
		str = fmt.Sprintf("%d. %s", d.displayIdx, text)
	}
	return str
}

func (d *dropdownItem) FilterValue() string {
	return d.comp.GetValue()
}

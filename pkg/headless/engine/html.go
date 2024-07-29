package engine

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"strings"
)

type UpdateHTMLFunc func(sel *goquery.Selection) error

func UpdateHTML(doc *goquery.Selection, sel DOMQuery, u UpdateHTMLFunc) error {
	out := sel.Query(doc)
	for _, v := range out {
		err := u(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateHTMLString(html string, sel DOMQuery, u UpdateHTMLFunc) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	return doc.Selection, UpdateHTML(doc.Selection, sel, u)
}

// DOMDataWriter represents anything that writes data into the DOM.
// Something like a mouse click, is not a DOMDataWriter.
type DOMDataWriter interface {
	GetQueries() []DOMQuery
}

// DOMQuery returns a list of selected HTML nodes given a query.
type DOMQuery interface {
	fmt.Stringer
	Query(doc *goquery.Selection) []*goquery.Selection
}

func NewEachSelector(c *v1beta1.EachSelector) *EachSelector {
	return &EachSelector{
		EachSelector: c,
	}
}

var _ DOMQuery = &EachSelector{}

type EachSelector struct {
	*v1beta1.EachSelector
}

func (e *EachSelector) Query(doc *goquery.Selection) []*goquery.Selection {
	sel := doc.Find(e.GetListSelector())
	out := make([]*goquery.Selection, 0, len(sel.Nodes))
	sel.Each(func(_ int, selection *goquery.Selection) {
		out = append(out, selection.Find(e.GetIterator()))
	})

	return out
}

var _ DOMQuery = StringSelector("")
var _ fmt.Stringer = StringSelector("")

type StringSelector string

func (s StringSelector) String() string {
	return string(s)
}

func (s StringSelector) Query(doc *goquery.Selection) []*goquery.Selection {
	return []*goquery.Selection{doc.Find(string(s))}
}

func FromChromeDev(node *cdp.Node) *v1beta1.HTMLNodeTemplateData {
	a := &v1beta1.HTMLNodeTemplateData{
		Attrs:    map[string]*v1beta1.HTMLNodeTemplateData_Attribute{},
		Element:  node.LocalName,
		Children: make([]*v1beta1.HTMLNodeTemplateData, 0, node.ChildNodeCount),
	}

	for i := 0; i < len(node.Attributes); i += 2 {
		name := node.Attributes[i]
		val := node.Attributes[i+1]
		a.Attrs[name] = &v1beta1.HTMLNodeTemplateData_Attribute{Val: val}
	}

	for _, v := range node.Children {
		if v.NodeType == cdp.NodeTypeText {
			a.Text = ptr.PtrOrNil(v.NodeValue)
		} else {
			a.Children = append(a.Children, FromChromeDev(v))
		}
	}

	return a
}

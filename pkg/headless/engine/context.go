package engine

import (
	"bytes"
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/flosch/pongo2/v6"
	"github.com/goccy/go-json"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/utils/collection"
	"golang.org/x/net/html"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Context struct {
	context.Context
	TemplateData *TemplateData
	Stream       renderer.Stream
	Evaluator    Evaluator
	Browser      Browser

	screenshotBuffer []byte
	opt              ContextOpts
}

type Browser interface {
	RefreshPage(ctx context.Context, content *string) error
	URL(ctx context.Context) (string, error)
	Exists(ctx context.Context, sel string) bool
	Navigate(ctx context.Context, addr string) error
	WriteInput(ctx context.Context, selector, content string) error
	Content(ctx context.Context, selector string) ([]*v1beta1.HTMLNodeTemplateData, error)
}

type ContextOpts struct {
	ScreenshotDir string
	// Where to save all the HTML files prior to every step run.
	HTMLDir string
}

func NewContext(parent context.Context, str renderer.Stream, eval Evaluator, o ContextOpts) (*Context, error) {
	out := &Context{
		Context: parent,
		TemplateData: &TemplateData{
			ScriptTemplateData: &v1beta1.ScriptTemplateData{
				Step:    make(map[string]string),
				Browser: &v1beta1.BrowserTemplateData{},
			},
		},
		Evaluator:        eval,
		Stream:           str,
		Browser:          NewBrowser(),
		opt:              o,
		screenshotBuffer: make([]byte, 10000),
	}

	return out, nil
}

func (c *Context) WithTimeout(to time.Duration) (*Context, context.CancelFunc) {
	if to == 0 {
		return c, func() {
		}
	}

	ctx, cancel := context.WithTimeout(c.Context, to)
	return c.WithContext(ctx), cancel
}

func (c *Context) WithContext(ctx context.Context) *Context {
	cp := c.Copy()
	cp.Context = ctx
	return cp
}

func (c *Context) Copy() *Context {
	return &Context{
		Context:          c.Context,
		TemplateData:     c.TemplateData,
		Stream:           c.Stream,
		Evaluator:        c.Evaluator,
		Browser:          c.Browser,
		screenshotBuffer: c.screenshotBuffer,
		opt:              c.opt,
	}
}

func (c *Context) RefreshPageState() error {
	err := c.Browser.RefreshPage(c.Context, &c.TemplateData.Browser.Content)
	if err != nil {
		return err
	}

	c.TemplateData.Browser.Url, err = c.Browser.URL(c.Context)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) SaveHTML(_ context.Context, fp, content string) error {
	err := os.WriteFile(fp, []byte(content), 0666)
	if err != nil {
		return errors.Join(except.NewFailed("failed to save HTML for step %s", fp), err)
	}

	return nil
}

func (c *Context) Screenshot(ctx context.Context, name string) (string, error) {
	if c.opt.ScreenshotDir == "" {
		return "", nil
	}

	for i, v := range c.screenshotBuffer {
		if v != 0 {
			c.screenshotBuffer[i] = 0
		}
	}

	err := chromedp.CaptureScreenshot(&c.screenshotBuffer).Do(ctx)
	if err != nil {
		return "", err
	}

	fp := filepath.Join(c.opt.ScreenshotDir, name+".png")
	err = os.WriteFile(fp, c.screenshotBuffer, 0666)
	if err != nil {
		return "", errors.Join(except.NewFailed("failed to save screenshot for step %s", name), err)
	}

	return fp, nil
}

func init() {
	_ = pongo2.RegisterFilter("html_select", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		query := param.String()
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(in.String()))
		if err != nil {
			return pongo2.AsValue(""), nil
		}

		s := doc.Find(query)

		var buf bytes.Buffer
		if len(s.Nodes) > 0 {
			for c := s.Nodes[0]; c != nil; c = c.NextSibling {
				err = html.Render(&buf, c)
				if err != nil {
					return pongo2.AsValue(""), nil
				}
			}
		}

		return pongo2.AsSafeValue(buf.String()), nil
	})
}

func CompileTemplate(expr string) (*TemplateRenderer, error) {
	tmp, err := pongo2.FromString(expr)
	if err != nil {
		return nil, err
	}
	return &TemplateRenderer{template: tmp}, nil
}

type TemplateRenderer struct {
	template *pongo2.Template
}

func (t *TemplateRenderer) Render(d *TemplateData) string {
	if t == nil {
		return ""
	}
	data := d.Copy()

	steps := data.Step
	data.Step = nil
	val := JSONVal(data)
	stepVal := pongo2.Context{}
	val["step"] = stepVal
	for k, v := range steps {
		addField(strings.Split(k, "."), v, stepVal)
	}

	v, _ := t.template.Execute(val)

	return v
}

func NodeToTemplate(node *html.Node) *v1beta1.HTMLNodeTemplateData {
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

type TemplateData struct {
	*v1beta1.ScriptTemplateData
}

func (t *TemplateData) Copy() *TemplateData {
	return &TemplateData{
		ScriptTemplateData: proto.Clone(t.ScriptTemplateData).(*v1beta1.ScriptTemplateData),
	}
}

func (t *TemplateData) ForEach(key string, f func(r *TemplateData)) {
	kp := strings.Split(key, ".")
	fe := map[string]string{}
	for k, v := range t.Step {
		spl := strings.Split(k, ".")
		if field, matched := collection.Rel(kp, spl); matched {
			fe[strings.Join(field, ".")] = v
		}
	}

	cl := t.Copy()
	cl.Each = fe
	f(cl)
}

func (t *TemplateData) GetStepVal(id string) string {
	if t.GetStep() == nil {
		return ""
	}

	return t.Step[id]
}

func (t *TemplateData) SetStepVal(id, val string) {
	if t.GetStep() == nil {
		t.Step = map[string]string{}
	}
	t.Step[id] = val
}

func (t *TemplateData) Merge(m ...*TemplateData) *TemplateData {
	out := t
	for _, v := range m {
		if v == nil {
			continue
		}
		proto.Merge(out, v.ScriptTemplateData)
	}

	return out
}

var (
	ProtoMarshaller = &protojson.MarshalOptions{
		AllowPartial: true,
	}

	ProtoUnmarshaller = &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
)

func ProtoJSONVal(msg proto.Message) pongo2.Context {
	b, _ := ProtoMarshaller.Marshal(msg)
	if len(b) == 0 {
		return nil
	}

	m := map[string]any{}
	_ = json.Unmarshal(b, &m)
	return m
}

func JSONVal(o any) pongo2.Context {
	b, _ := json.Marshal(o)
	if len(b) == 0 {
		return nil
	}

	m := map[string]any{}
	_ = json.Unmarshal(b, &m)
	return m
}

func MergeTemplateContexts(c ...pongo2.Context) pongo2.Context {
	out := pongo2.Context{}
	for _, v := range c {
		for k, val := range v {
			out[k] = val
		}
	}
	return out
}

func IsTruthy(a any) bool {
	if a == nil {
		return false
	}

	switch t := a.(type) {
	case string:
		t = strings.TrimSpace(t)
		return t != "" && strings.ToLower(t) != "false"
	case bool:
		return t
	case float32:
		return t != 0
	case float64:
		return t != 0
	case int:
		return t != 0
	case int8:
		return t != 0
	case int16:
		return t != 0
	case int32:
		return t != 0
	case int64:
		return t != 0
	case uint:
		return t != 0
	case uint8:
		return t != 0
	case uint16:
		return t != 0
	case uint32:
		return t != 0
	case uint64:
		return t != 0
	}
	return false
}

func addField(p []string, val string, c pongo2.Context) {
	switch len(p) {
	case 0:
		return
	case 1:
		c[p[0]] = val
		return
	default:
	}

	field := p[0]

	var temp pongo2.Context
	q, ok := c[field]
	if !ok {
		temp = pongo2.Context{}
		c[field] = temp
	} else {
		temp = q.(pongo2.Context)
	}

	addField(p[1:], val, temp)
}

func NewBrowser() Browser {
	return &browser{}
}

type browser struct {
}

func (p *browser) Content(ctx context.Context, selector string) ([]*v1beta1.HTMLNodeTemplateData, error) {
	out := make([]*v1beta1.HTMLNodeTemplateData, 0)
	var nodes []*cdp.Node
	err := chromedp.Nodes(selector, &nodes).Do(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range nodes {
		a := &v1beta1.HTMLNodeTemplateData{
			Data:    v.NodeValue,
			Attrs:   map[string]*v1beta1.HTMLNodeTemplateData_Attribute{},
			Element: v.LocalName,
		}

		for i := 0; i < len(v.Attributes); i += 2 {
			name := v.Attributes[i]
			val := v.Attributes[i+1]
			a.Attrs[name] = &v1beta1.HTMLNodeTemplateData_Attribute{Val: val}
		}

		for _, v := range v.Children {
			if v.NodeType == cdp.NodeTypeText {
				a.Data = v.NodeValue
			}
		}

		out = append(out, a)
	}

	return out, err
}

func (p *browser) WriteInput(ctx context.Context, selector, content string) error {
	return chromedp.SendKeys(selector, content).Do(ctx)
}

func (p *browser) Navigate(ctx context.Context, addr string) error {
	return chromedp.Navigate(addr).Do(ctx)
}

func (p *browser) Exists(ctx context.Context, sel string) (exists bool) {
	_ = chromedp.QueryAfter(sel, func(ctx context.Context, id runtime.ExecutionContextID, node ...*cdp.Node) error {
		exists = len(node) > 0
		return nil
	}, chromedp.RetryInterval(0)).Do(ctx)
	return
}

func (p *browser) URL(ctx context.Context) (s string, err error) {
	err = chromedp.Location(&s).Do(ctx)
	return
}

func (p *browser) RefreshPage(ctx context.Context, content *string) error {
	err := chromedp.OuterHTML("html", content).Do(ctx)
	return err
}

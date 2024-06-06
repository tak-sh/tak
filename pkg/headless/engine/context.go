package engine

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
	"github.com/flosch/pongo2/v6"
	"github.com/goccy/go-json"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/renderer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"strings"
)

type Context struct {
	context.Context
	TemplateData *TemplateData
	Stream       renderer.Stream

	screenshotBuffer []byte
	opt              ContextOpts
}

type ContextOpts struct {
	ScreenshotDir string
	// Where to save all the HTML files prior to every step run.
	HTMLDir string
}

func NewContext(parent context.Context, str renderer.Stream, o ContextOpts) (*Context, error) {
	out := &Context{
		Context: parent,
		TemplateData: &TemplateData{
			ScriptTemplateData: &v1beta1.ScriptTemplateData{
				Step:    make(map[string]string),
				Browser: &v1beta1.BrowserTemplateData{},
			},
		},
		Stream:           str,
		opt:              o,
		screenshotBuffer: make([]byte, 10000),
	}

	return out, nil
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

type TemplateData struct {
	*v1beta1.ScriptTemplateData
}

func (t *TemplateData) Render(v string) string {
	tmp, err := pongo2.FromString(v)
	if err != nil {
		return ""
	}

	out, err := tmp.Execute(JSONVal(t.ScriptTemplateData))
	if err != nil {
		return ""
	}
	return out
}

func (t *TemplateData) Merge(m ...*TemplateData) *TemplateData {
	out := proto.Clone(t.ScriptTemplateData).(*v1beta1.ScriptTemplateData)
	for _, v := range m {
		proto.Merge(out, v.ScriptTemplateData)
	}

	return &TemplateData{out}
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

package headless

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
	"github.com/flosch/pongo2/v6"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/renderer"
	"os"
	"path/filepath"
	"strings"
)

type Context struct {
	context.Context
	Store  Store
	Stream renderer.Stream

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
		Context:          parent,
		Store:            map[string]any{},
		Stream:           str,
		opt:              o,
		screenshotBuffer: make([]byte, 10000),
	}

	return out, nil
}

func (c *Context) SaveHTML(_ context.Context, name, content string) error {
	if c.opt.HTMLDir == "" {
		return nil
	}

	fp := filepath.Join(c.opt.HTMLDir, name+".html")
	err := os.WriteFile(fp, []byte(content), 0666)
	if err != nil {
		return errors.Join(except.NewFailed("failed to save HTML for step %s", name), err)
	}

	return nil
}

func (c *Context) Screenshot(ctx context.Context, name string) error {
	if c.opt.ScreenshotDir == "" {
		return nil
	}

	for i, v := range c.screenshotBuffer {
		if v != 0 {
			c.screenshotBuffer[i] = 0
		}
	}

	err := chromedp.CaptureScreenshot(&c.screenshotBuffer).Do(ctx)
	if err != nil {
		return err
	}

	fp := filepath.Join(c.opt.ScreenshotDir, name+".png")
	err = os.WriteFile(fp, c.screenshotBuffer, 0666)
	if err != nil {
		return errors.Join(except.NewFailed("failed to save screenshot for step %s", name), err)
	}

	return nil
}

type Store pongo2.Context

func (s Store) Set(id string, val any) {
	s[id] = val
}

func (s Store) Get(id string) any {
	return s[id]
}

func (s Store) Render(v string) string {
	tmp, err := pongo2.FromString(v)
	if err != nil {
		return ""
	}

	out, _ := tmp.Execute(pongo2.Context(s))
	return out
}

func (s Store) Merge(m ...Store) Store {
	out := Store{}
	for k, v := range s {
		out[k] = v
	}

	for _, v := range m {
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

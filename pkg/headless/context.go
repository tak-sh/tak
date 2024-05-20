package headless

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
	"github.com/eddieowens/opts"
	"github.com/flosch/pongo2/v6"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/prompt"
	"os"
	"path/filepath"
)

type Context struct {
	context.Context
	Store  Store
	Stream prompt.Stream

	screenshotBuffer []byte
	opt              ContextOpts
}

type ContextOpts struct {
	ScreenshotDir string
}

func WithScreenshotsDir(dir string) opts.Opt[ContextOpts] {
	return func(c *ContextOpts) {
		c.ScreenshotDir = dir
	}
}

func (c ContextOpts) DefaultOptions() ContextOpts {
	return ContextOpts{}
}

func NewContext(parent context.Context, str prompt.Stream, o ...opts.Opt[ContextOpts]) (*Context, error) {
	out := &Context{
		Context:          parent,
		Store:            map[string]any{},
		Stream:           str,
		opt:              opts.DefaultApply(o...),
		screenshotBuffer: make([]byte, 10000),
	}

	return out, nil
}

func (c *Context) RenderTemplate(template string) string {
	tmp, err := pongo2.FromString(template)
	if err != nil {
		return ""
	}

	out, _ := tmp.Execute(pongo2.Context(c.Store))
	return out
}

func (c *Context) Screenshot(ctx context.Context, id string) error {
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

	fp := filepath.Join(c.opt.ScreenshotDir, id+".png")
	err = os.WriteFile(fp, c.screenshotBuffer, 0666)
	if err != nil {
		return errors.Join(except.NewFailed("failed to save screenshot for step %s", id), err)
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

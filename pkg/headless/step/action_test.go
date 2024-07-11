package step

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"testing/fstest"
)

type ActionTestSuite struct {
	suite.Suite
}

func (a *ActionTestSuite) TestAct() {
	type test struct {
		Given           *v1beta1.Action
		Id              string
		Post            func(desc string)
		ExpectedErr     string
		ExpectedCompErr string
		Ctx             *engine.Context
	}

	tests := map[string]test{
		"store works": func() test {
			c, _ := engine.NewContext(context.Background(), nil, nil, engine.ContextOpts{})

			return test{
				Given: &v1beta1.Action{
					Store: &v1beta1.Action_Store{
						KeyVals: []*v1beta1.KeyVal{
							{Key: "derp.flerp", Value: "1"},
						},
					},
				},
				Id: "derpflorp",
				Post: func(desc string) {
					o := c.TemplateData.GetStepVal("derpflorp.derp.flerp")
					a.Equal("1", o, desc)
				},
				Ctx: c,
			}
		}(),
	}

	for desc, t := range tests {
		given, err := New(t.Id, t.Given)
		if err != nil {
			if t.ExpectedCompErr != "" {
				a.EqualError(err, t.ExpectedCompErr, desc)
			} else {
				a.NoError(err)
			}
			continue
		}

		err = given.Eval(t.Ctx, 0)
		if err != nil {
			if t.ExpectedErr != "" {
				a.EqualError(err, t.ExpectedErr, desc)
			} else {
				a.NoError(err)
			}
			continue
		}

		if t.Post != nil {
			t.Post(desc)
		}
	}
}

func (a *ActionTestSuite) TestActWithBrowser() {
	type test struct {
		GivenCtx           *engine.Context
		GivenID            string
		GivenHTML          fstest.MapFS
		Given              *v1beta1.Action
		URL                string
		Validate           func(desc string)
		ExpectedCompileErr string
		ExpectedErr        string
	}

	tests := map[string]test{
		"asd": func() test {
			c, _ := engine.NewContext(context.Background(), nil, nil, engine.ContextOpts{})
			return test{
				Validate: func(desc string) {
					a.Equal(map[string]string{
						"id.One":   "1",
						"id.Two":   "2",
						"id.Three": "3",
					}, c.TemplateData.Step, desc)
				},
				GivenCtx: c,
				GivenID:  "id",
				Given: &v1beta1.Action{
					ForEachElement: &v1beta1.Action_ForEachElement{
						Selector: "ul[class='list'] > li",
						Actions: []*v1beta1.Action{
							{Store: &v1beta1.Action_Store{
								KeyVals: []*v1beta1.KeyVal{
									{Key: "{{element.data}}", Value: "{{element.attrs.class.val}}"},
								},
							}},
						},
					},
				},
				GivenHTML: fstest.MapFS{
					"index.html": {
						Data: []byte(`
    <ul class="list">
        <li class="1">One</li>
        <li class="2">Two</li>
        <li class="3">Three</li>
    </ul>
`),
					},
				},
			}
		}(),
	}

	for desc, t := range tests {
		runner, err := NewServer(t.GivenCtx, t.GivenHTML, t.GivenID, t.Given)
		if err != nil {
			if t.ExpectedCompileErr != "" {
				a.EqualError(err, t.ExpectedCompileErr, desc)
			} else {
				a.NoError(err, desc)
			}
			continue
		}

		err = runner(t.URL)
		if err != nil {
			if t.ExpectedErr != "" {
				a.EqualError(err, t.ExpectedErr, desc)
			} else {
				a.NoError(err, desc)
			}
			continue
		}

		t.Validate(desc)
	}
}

func NewServer(c *engine.Context, h fs.FS, id string, acts ...*v1beta1.Action) (func(u string) error, error) {
	s := httptest.NewServer(http.FileServerFS(h))
	ctx := context.Background()
	ctx, cancel := chromedp.NewContext(ctx)

	compiled := make([]Action, 0, len(acts))
	for _, v := range acts {
		a, err := New(id, v)
		if err != nil {
			return nil, err
		}

		compiled = append(compiled, a)
	}

	return func(pa string) error {
		defer func() {
			s.Close()
			cancel()
		}()
		return chromedp.Run(ctx,
			chromedp.Navigate(path.Join(s.URL, pa)),
			chromedp.ActionFunc(func(ctx context.Context) error {
				c.Context = ctx
				for _, v := range compiled {
					err := v.Eval(c, 0)
					if err != nil {
						return err
					}
				}
				return nil
			}),
		)
	}, nil
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(ActionTestSuite))
}

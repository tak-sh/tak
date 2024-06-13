package engine

import (
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"testing"
)

type ContextTestSuite struct {
	suite.Suite
}

func (c *ContextTestSuite) TestRenderTemplateData() {
	type test struct {
		Data     *v1beta1.ScriptTemplateData
		Given    string
		Expected string
	}

	tests := map[string]test{
		"renders dots properly": {
			Data: &v1beta1.ScriptTemplateData{
				Step: map[string]string{
					"mfa.2fa": "derp",
				},
			},
			Given:    "{{step.mfa.2fa}}",
			Expected: "derp",
		},
		"renders single level properly": {
			Data: &v1beta1.ScriptTemplateData{
				Step: map[string]string{
					"mfa": "derp",
				},
			},
			Given:    "{{step.mfa}}",
			Expected: "derp",
		},
		"renders empty step": {
			Data:  &v1beta1.ScriptTemplateData{},
			Given: "{{step.mfa}}",
		},
		"renders selector": {
			Data: &v1beta1.ScriptTemplateData{
				Browser: &v1beta1.BrowserTemplateData{
					Content: `<!doctype html><html lang="en"><div><input id="123"/></div><body></body></html>`,
				},
			},
			Expected: `<input id="123"/>`,
			Given:    `{{browser.content|html_select:"input[id='123']"}}`,
		},
		"renders empty selector": {
			Data: &v1beta1.ScriptTemplateData{
				Browser: &v1beta1.BrowserTemplateData{
					Content: `<!doctype html><html lang="en"><div><input id="123"/></div><body></body></html>`,
				},
			},
			Given: `{{browser.content|html_select:"input[id='456']"}}`,
		},
	}

	for desc, v := range tests {
		temp, err := CompileTemplate(v.Given)
		if !c.NoError(err, desc) {
			continue
		}

		actual := temp.Render(&TemplateData{ScriptTemplateData: v.Data})
		c.Equal(v.Expected, actual, desc)
	}
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}

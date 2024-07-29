package component

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"testing"
	"time"
)

type DropdownTestSuite struct {
	TestSuite
}

func (d *DropdownTestSuite) TestRender() {
	type test struct {
		GivenHTML       string
		TemplateData    *engine.TemplateData
		ExpectedOptions []*v1beta1.Component_Dropdown_Option
		Given           *v1beta1.Component_Dropdown
	}

	tests := map[string]test{
		"render from HTML": {
			GivenHTML: `<div class="col-xs-12 row"> <div class="jpui styledselect show" id="simplerAuth-dropdownoptions-styledselect">  <span class="wrap right" id="span-header-simplerAuth-dropdownoptions-styledselect"><input type="button" id="header-simplerAuth-dropdownoptions-styledselect" name="contact" form="" class="jpui input header focus-on-header wrap right text-float-left" aria-haspopup="true" aria-disabled="false" aria-expanded="true" aria-label=" Tell us how: Choose one" value="Choose one">  <button class="jpui input-icon icon text-overflow" id="iconButton-simplerAuth-dropdownoptions-styledselect" type="button" aria-hidden="true" tabindex="-1"><span class="jpui expanddown icon input-icon hasError"></span></button></span> <div class="list-container open"><ul class="list" id="ul-list-container-simplerAuth-dropdownoptions-styledselect" role="listbox"><li role="presentation"><a class="option js-option groupLabelContainer STYLED_SELECT" id="container-0-simplerAuth-dropdownoptions-styledselect" href="javascript:void(0);" rel="" role="option" aria-setsize="7" aria-posinset="1" tabindex="0" aria-disabled="true"><span class="groupLabelText primary" id="container-primary-0-simplerAuth-dropdownoptions-styledselect"> TEXT ME</span><span class="secondary" id="container-secondary-0-simplerAuth-dropdownoptions-styledselect"> </span><span class="util accessible-text" id="option-accessible-0-simplerAuth-dropdownoptions-styledselect">   </span></a></li><li role="presentation"><a class="option js-option STYLED_SELECT" id="container-1-simplerAuth-dropdownoptions-styledselect" href="javascript:void(0);" rel="S492295510" role="option" aria-setsize="7" aria-posinset="2" tabindex="0"><span class="primary groupingName" id="container-primary-1-simplerAuth-dropdownoptions-styledselect">xxx-xxx-5621</span><span class="secondary" id="container-secondary-1-simplerAuth-dropdownoptions-styledselect"> </span><span class="util accessible-text" id="option-accessible-1-simplerAuth-dropdownoptions-styledselect">   </span></a></li><li role="presentation"><a class="option js-option lastGroupItem STYLED_SELECT" id="container-2-simplerAuth-dropdownoptions-styledselect" href="javascript:void(0);" rel="S634050098" role="option" aria-setsize="7" aria-posinset="3" tabindex="0"><span class="primary groupingName" id="container-primary-2-simplerAuth-dropdownoptions-styledselect">xxx-xxx-5622</span><span class="secondary" id="container-secondary-2-simplerAuth-dropdownoptions-styledselect"> </span><span class="util accessible-text" id="option-accessible-2-simplerAuth-dropdownoptions-styledselect">   </span></a></li><li role="presentation"><a class="option js-option groupLabelContainer STYLED_SELECT" id="container-3-simplerAuth-dropdownoptions-styledselect" href="javascript:void(0);" rel="" role="option" aria-setsize="7" aria-posinset="4" tabindex="0" aria-disabled="true"><span class="groupLabelText primary" id="container-primary-3-simplerAuth-dropdownoptions-styledselect"> CALL ME</span><span class="secondary" id="container-secondary-3-simplerAuth-dropdownoptions-styledselect"> </span><span class="util accessible-text" id="option-accessible-3-simplerAuth-dropdownoptions-styledselect">   </span></a></li><li role="presentation"><a class="option js-option STYLED_SELECT" id="container-4-simplerAuth-dropdownoptions-styledselect" href="javascript:void(0);" rel="V492295510" role="option" aria-setsize="7" aria-posinset="5" tabindex="0"><span class="primary groupingName" id="container-primary-4-simplerAuth-dropdownoptions-styledselect">xxx-xxx-5621</span><span class="secondary" id="container-secondary-4-simplerAuth-dropdownoptions-styledselect"> </span><span class="util accessible-text" id="option-accessible-4-simplerAuth-dropdownoptions-styledselect">   </span></a></li><li role="presentation"><a class="option js-option styledSelectSeparator lastGroupItem STYLED_SELECT" id="container-5-simplerAuth-dropdownoptions-styledselect" href="javascript:void(0);" rel="V634050098" role="option" aria-setsize="7" aria-posinset="6" tabindex="0"><span class="primary groupingName" id="container-primary-5-simplerAuth-dropdownoptions-styledselect">xxx-xxx-5622</span><span class="secondary" id="container-secondary-5-simplerAuth-dropdownoptions-styledselect"> </span><span class="util accessible-text" id="option-accessible-5-simplerAuth-dropdownoptions-styledselect">   </span></a></li><li role="presentation"><a class="option js-option STYLED_SELECT" id="container-6-simplerAuth-dropdownoptions-styledselect" href="javascript:void(0);" rel="Call" role="option" aria-setsize="7" aria-posinset="7" tabindex="0"><span class="primary" id="container-primary-6-simplerAuth-dropdownoptions-styledselect">Call us - 1-877-242-7372</span><span class="secondary" id="container-secondary-6-simplerAuth-dropdownoptions-styledselect"> </span><span class="util accessible-text" id="option-accessible-6-simplerAuth-dropdownoptions-styledselect">   </span></a></li></ul></div>  </div></div>`,
			Given: &v1beta1.Component_Dropdown{
				From: &v1beta1.Component_Dropdown_FromSpec{
					Selector: &v1beta1.EachSelector{
						ListSelector: "#ul-list-container-simplerAuth-dropdownoptions-styledselect > li",
						Iterator:     "span[id*='container-primary']",
					},
					Mapper: &v1beta1.Component_Dropdown_Option{
						Value: "{{ element.attrs.id.val }}",
						Text:  ptr.Ptr("{{ element.text }}"),
					},
				},
				Merge: []*v1beta1.Component_Dropdown_OptionMerge{
					{
						If: ptr.Ptr("{{ 'Call us' in option.text }}"),
						Option: &v1beta1.Component_Dropdown_Option{
							Hidden: ptr.Ptr(true),
						},
					},
					{
						If: ptr.Ptr(`{{ not ('xxx-' in option.text) }}`),
						Option: &v1beta1.Component_Dropdown_Option{
							Disabled: ptr.Ptr(true),
						},
					},
				},
			},
			ExpectedOptions: []*v1beta1.Component_Dropdown_Option{
				{Text: ptr.Ptr(" TEXT ME"), Value: "container-primary-0-simplerAuth-dropdownoptions-styledselect", Disabled: ptr.Ptr(true)},
				{Text: ptr.Ptr("xxx-xxx-5621"), Value: "container-primary-1-simplerAuth-dropdownoptions-styledselect"},
				{Text: ptr.Ptr("xxx-xxx-5622"), Value: "container-primary-2-simplerAuth-dropdownoptions-styledselect"},
				{Text: ptr.Ptr(" CALL ME"), Value: "container-primary-3-simplerAuth-dropdownoptions-styledselect", Disabled: ptr.Ptr(true)},
				{Text: ptr.Ptr("xxx-xxx-5621"), Value: "container-primary-4-simplerAuth-dropdownoptions-styledselect"},
				{Text: ptr.Ptr("xxx-xxx-5622"), Value: "container-primary-5-simplerAuth-dropdownoptions-styledselect"},
			},
		},
	}

	for desc, v := range tests {
		dd, _ := NewDropdown(v.Given)

		expectedItems := make([]*dropdownItem, 0, len(v.ExpectedOptions))
		for i := range v.ExpectedOptions {
			v := v.ExpectedOptions[i]
			expectedItems = append(expectedItems, &dropdownItem{comp: v})
		}

		c, _ := engine.NewContext(context.Background(), nil, engine.NewEvaluator(engine.NewEventQueue(), 1*time.Second), engine.ContextOpts{})
		c.TemplateData.Browser.Content = v.GivenHTML
		mod := dd.Render(c, &Props{})
		actual := mod.(*DropdownModel)
		d.EqualDropdownItems(expectedItems, actual.List.Items(), desc)
	}

}

func TestDropdownTestSuite(t *testing.T) {
	suite.Run(t, new(DropdownTestSuite))
}

package step

import (
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/component"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"github.com/tak-sh/tak/pkg/validate"
)

var _ Action = &PromptAction{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_PromptUser] = &PromptAction{}

type PromptAction struct {
	prompt *v1beta1.Action_PromptUser
	Prompt *Prompt
	ID     string
}

func NewPromptAction(id string, p *v1beta1.Action_PromptUser) (*PromptAction, error) {
	out := &PromptAction{
		ID:     id,
		prompt: p,
	}

	var err error
	out.Prompt, err = NewPrompt(p.GetPrompt())
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (p *PromptAction) Validate() error {
	err := p.Prompt.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (p *PromptAction) String() string {
	prmpt := p.prompt.GetPrompt().GetTitle()
	if p.prompt.GetPrompt().Description != nil {
		prmpt = p.prompt.GetPrompt().GetDescription()
	}
	return fmt.Sprintf("asking the user %s", prmpt)
}

func (p *PromptAction) Act(ctx *engine.Context) error {
	model := p.Prompt.Component.Render(ctx, &component.Props{
		ID:          p.GetID(),
		Title:       component.TitleStyle.Render(p.prompt.GetPrompt().GetTitle()),
		Description: component.DescriptionStyle.Render(p.prompt.GetPrompt().GetDescription()),
	})
	if model == nil {
		return nil
	}

	v, err := ctx.Stream.Render(ctx, model)
	if err != nil {
		return err
	}

	ctx.TemplateData.SetStepVal(p.ID, GetValueString(v.Value))

	return nil
}

func (p *PromptAction) GetID() string {
	return p.ID
}

func (p *PromptAction) ToProto() *v1beta1.Action_PromptUser {
	return p.prompt
}

func NewPrompt(p *v1beta1.Prompt) (*Prompt, error) {
	cmp, err := component.New(p.GetComponent())
	if err != nil {
		return nil, err
	}

	out := &Prompt{
		prompt:    p,
		Component: cmp,
	}

	return out, nil
}

var _ grpcutils.ProtoWrapper[*v1beta1.Prompt] = &Prompt{}
var _ validate.Validator = &Prompt{}

type Prompt struct {
	prompt    *v1beta1.Prompt
	Component component.Component
}

func (p *Prompt) Validate() error {
	if p.prompt.GetTitle() == "" {
		return except.NewInvalid("a title is required")
	}

	err := p.Component.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (p *Prompt) ToProto() *v1beta1.Prompt {
	return p.prompt
}

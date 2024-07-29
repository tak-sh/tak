package provider

import (
	"errors"
	"github.com/tak-sh/tak/generated/go/api/provider/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
)

type AccountMapping struct {
	*v1beta1.ListAccountsSpec_AccountMapping

	NameCompiled *engine.TemplateRenderer
	TypeCompiled *engine.TemplateRenderer
}

func newAccountMapping(act *v1beta1.ListAccountsSpec_AccountMapping) (*AccountMapping, error) {
	if act == nil {
		return nil, except.NewInvalid("missing")
	}

	out := &AccountMapping{
		ListAccountsSpec_AccountMapping: act,
	}

	var err error
	out.NameCompiled, err = engine.CompileTemplate(act.Name)
	if err != nil {
		return nil, errors.Join(except.NewInvalid("name field template is invalid"), err)
	}

	out.TypeCompiled, err = engine.CompileTemplate(act.Type)
	if err != nil {
		return nil, errors.Join(except.NewInvalid("type field template is invalid"), err)
	}

	return out, nil
}

func (a *AccountMapping) Render(d *engine.TemplateData) *v1beta1.Account {
	return &v1beta1.Account{
		Name: a.NameCompiled.Render(d),
		Type: v1beta1.AccountType_Enum(v1beta1.AccountType_Enum_value[a.TypeCompiled.Render(d)]),
	}
}

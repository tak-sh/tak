package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"strings"
)

var _ Action = &StoreAction{}

type StoreAction struct {
	*v1beta1.Action_Store
	CompiledKeyVals []*KeyVal
	ID              string
}

type KeyVal struct {
	*v1beta1.KeyVal
	CompiledKey   *engine.TemplateRenderer
	CompiledValue *engine.TemplateRenderer
}

func NewStoreAction(id string, s *v1beta1.Action_Store) (*StoreAction, error) {
	out := &StoreAction{
		ID:              id,
		Action_Store:    s,
		CompiledKeyVals: make([]*KeyVal, 0, len(s.GetKeyVals())),
	}

	for i, v := range s.GetKeyVals() {
		kv, err := NewKeyVal(v)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("key val %d", i), err)
		}
		out.CompiledKeyVals = append(out.CompiledKeyVals, kv)
	}

	return out, nil
}

func NewKeyVal(kv *v1beta1.KeyVal) (*KeyVal, error) {
	out := &KeyVal{
		KeyVal: kv,
	}

	var err error
	out.CompiledKey, err = engine.CompileTemplate(kv.GetKey())
	if err != nil {
		return nil, errors.Join(except.NewInvalid("key %s", kv.GetKey()), err)
	}

	out.CompiledValue, err = engine.CompileTemplate(kv.GetValue())
	if err != nil {
		return nil, errors.Join(except.NewInvalid("value %s", kv.GetValue()), err)
	}

	return out, nil
}

func (s *StoreAction) String() string {
	strs := make([]string, 0, len(s.GetKeyVals()))
	return "storing " + strings.Join(strs, ", ")
}

func (s *StoreAction) Validate() error {
	return nil
}

func (s *StoreAction) Act(ctx *engine.Context) error {
	for _, kv := range s.CompiledKeyVals {
		keys := strings.Join([]string{s.ID, kv.CompiledKey.Render(ctx.TemplateData)}, ".")
		ctx.TemplateData.SetStepVal(keys, kv.CompiledValue.Render(ctx.TemplateData))
	}
	return nil
}

func (s *StoreAction) GetID() string {
	return s.ID
}

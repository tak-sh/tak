package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

var _ Action = &StoreAction{}

type StoreAction struct {
	*v1beta1.Action_Store
	CompiledKeyVals []*KeyVal
	ID              string
}

func (s *StoreAction) Message() proto.Message {
	return s.Action_Store
}

func (s *StoreAction) GetId() string {
	return s.ID
}

func (s *StoreAction) Eval(c *engine.Context, to time.Duration) error {
	c, cancel := c.WithTimeout(to)
	defer cancel()

	for _, kv := range s.CompiledKeyVals {
		keys := strings.Join([]string{s.ID, kv.CompiledKey.Render(c.TemplateData)}, ".")
		c.TemplateData.SetStepVal(keys, kv.CompiledValue.Render(c.TemplateData))
	}
	return nil
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

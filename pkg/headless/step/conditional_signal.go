package step

import (
	"errors"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
)

var _ grpcutils.ProtoWrapper[*v1beta1.ConditionalSignal] = &ConditionalSignal{}
var _ Signaller = &ConditionalSignal{}

func NewConditionalSignal(s *v1beta1.ConditionalSignal) (*ConditionalSignal, error) {
	out := &ConditionalSignal{
		ConditionalSignal: s,
	}

	if s.If == "" {
		return nil, except.NewInvalid("'if' field is required for a signal")
	}

	if s.Signal == v1beta1.ConditionalSignal_unknown {
		return nil, except.NewInvalid("'signal' field is required for a signal")
	}

	if s.Signal == v1beta1.ConditionalSignal_error && s.GetMessage() == "" {
		return nil, except.NewInvalid("'message' field required for any error signals")
	}

	var err error
	out.Conditional, err = engine.CompileTemplate(s.GetIf())
	if err != nil {
		errors.Join(except.NewInvalid("signals 'if' field"), err)
	}

	return out, nil
}

type ConditionalSignal struct {
	*v1beta1.ConditionalSignal
	Conditional *engine.TemplateRenderer
}

func (s *ConditionalSignal) CheckSignal(st *engine.TemplateData) *v1beta1.ConditionalSignal {
	if s == nil {
		return nil
	}

	if engine.IsTruthy(s.Conditional.Render(st)) {
		return s.ConditionalSignal
	}

	return nil
}

func (s *ConditionalSignal) ToProto() *v1beta1.ConditionalSignal {
	if s == nil {
		return nil
	}

	return s.ConditionalSignal
}

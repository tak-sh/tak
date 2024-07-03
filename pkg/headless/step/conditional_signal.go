package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"strings"
)

var _ grpcutils.ProtoWrapper[*v1beta1.ConditionalSignal] = &ConditionalSignal{}
var _ engine.PathNode = &ConditionalSignal{}

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
		return nil, errors.Join(except.NewInvalid("signals 'if' field"), err)
	}

	return out, nil
}

type ConditionalSignal struct {
	*v1beta1.ConditionalSignal
	Conditional *engine.TemplateRenderer
}

func (s *ConditionalSignal) String() string {
	switch s.Signal {
	case v1beta1.ConditionalSignal_success:
		out := make([]string, 0, 2)
		out = append(out, "success")
		if s.GetMessage() != "" {
			out = append(out, s.GetMessage())
		}
		return strings.Join(out, ": ")
	case v1beta1.ConditionalSignal_error:
		return fmt.Sprintf("error: %s", s.GetMessage())
	default:
		return ""
	}
}

func (s *ConditionalSignal) GetId() string {
	return ""
}

func (s *ConditionalSignal) IsReady(st *engine.Context) bool {
	return engine.IsTruthy(s.Conditional.Render(st.TemplateData))
}

func (s *ConditionalSignal) ToProto() *v1beta1.ConditionalSignal {
	if s == nil {
		return nil
	}

	return s.ConditionalSignal
}

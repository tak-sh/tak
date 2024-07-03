package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"github.com/tak-sh/tak/pkg/validate"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewStep(s *v1beta1.Step) (*Step, error) {
	act, err := New(s.GetId(), s.GetAction())
	if err != nil {
		return nil, err
	}

	out := &Step{
		CompiledAction:     act,
		Step:               s,
		ConditionalSignals: make([]*ConditionalSignal, len(s.GetSignals())),
		cancelChan:         make(chan error, 1),
	}

	for i, v := range s.GetSignals() {
		out.ConditionalSignals[i], err = NewConditionalSignal(v)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("signal #%d", i), err)
		}
	}

	return out, nil
}

var idRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

func NewID(parent string, idx int, s *v1beta1.Step) (string, error) {
	newId := make([]string, 0, 2)
	newId = append(newId, parent)

	if s.Id != nil {
		id := *s.Id
		if !idRegex.MatchString(id) {
			return "", except.NewInvalid("%s is not a valid id. IDs must contain alphanumeric, '_', '-' characters", id)
		}
		newId = append(newId, id)
	} else {
		newId = append(newId, strconv.Itoa(idx))
	}
	return strings.Join(newId, "."), nil
}

var _ grpcutils.ProtoWrapper[*v1beta1.Step] = &Step{}
var _ validate.Validator = &Step{}
var _ engine.Instruction = &Step{}
var _ engine.PathNode = &Step{}

type Step struct {
	*v1beta1.Step
	CompiledAction     Action
	ConditionalSignals []*ConditionalSignal

	cancelChan chan error
}

func (s *Step) String() string {
	return s.CompiledAction.String()
}

func (s *Step) Cancel(err error) {
	if s == nil {
		return
	}
	if len(s.cancelChan) == 0 {
		s.cancelChan <- err
	}
}

func (s *Step) IsReady(c *engine.Context) bool {
	v, ok := s.CompiledAction.(engine.PathNode)
	if ok {
		return v.IsReady(c)
	}

	return true
}

func (s *Step) drainCancelChan() {
	for {
		select {
		case _, ok := <-s.cancelChan:
			if !ok {
				return
			}
		default:
			return
		}
	}
}

func (s *Step) Eval(c *engine.Context, to time.Duration) (err error) {
	s.drainCancelChan()
	out := RunActionAsync(c, s.CompiledAction, to)
	select {
	case err = <-out:
	case err = <-s.cancelChan:
	}
	return
}

func (s *Step) Validate() error {
	if s.GetAction().GetAsk() != nil && s.Id == nil {
		return except.NewInvalid("any step with a prompt must have an ID")
	}

	err := s.CompiledAction.Validate()
	if err != nil {
		if s.Id != nil {
			return errors.Join(fmt.Errorf("id %s", s.GetId()), err)
		}
		return err
	}

	return nil
}

func (s *Step) ToProto() *v1beta1.Step {
	return s.Step
}

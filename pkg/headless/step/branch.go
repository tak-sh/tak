package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"time"
)

func NewBranch(id string, b *v1beta1.Action_Branch) (*BranchAction, error) {
	out := &BranchAction{
		Action_Branch: b,
		ID:            id,
		CompiledSteps: make([]*Step, 0, len(b.Steps)),
	}

	if b.If == "" {
		return nil, except.NewInvalid("branches must have an 'if' statement")
	}

	var err error
	out.ShouldRunCond, err = engine.CompileTemplate(b.GetIf())
	if err != nil {
		return nil, errors.Join(except.NewInvalid("%s is not a valid if", b.GetIf()), err)
	}

	for i, v := range b.GetSteps() {
		stepId, err := NewID(id, i, v)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("branch %s, step %d", id, i), err)
		}
		v.Id = &stepId
		st, err := NewStep(v)
		if err != nil {
			return nil, err
		}

		out.CompiledSteps = append(out.CompiledSteps, st)
	}

	return out, nil
}

type Branches interface {
	Statements() []engine.Instruction
}

var _ Action = &BranchAction{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_Branch] = &BranchAction{}
var _ engine.PathNode = &BranchAction{}
var _ Branches = &BranchAction{}

type BranchAction struct {
	*v1beta1.Action_Branch
	ID            string
	CompiledSteps []*Step
	ShouldRunCond *engine.TemplateRenderer
}

func (b *BranchAction) Eval(c *engine.Context, _ time.Duration) error {
	if !engine.IsTruthy(b.ShouldRunCond.Render(c.TemplateData)) {
		return nil
	}

	for _, v := range b.CompiledSteps {
		handle := c.Evaluator.Eval(c, v.CompiledAction)
		<-handle.Done()
		if err := handle.Cause(); err != nil {
			return errors.Join(fmt.Errorf("step %s", v.GetId()), err)
		}
	}

	return nil
}

func (b *BranchAction) Statements() []engine.Instruction {
	return nil
}

func (b *BranchAction) GetId() string {
	return b.ID
}

func (b *BranchAction) ShouldBranch(c *engine.Context) bool {
	return engine.IsTruthy(b.ShouldRunCond.Render(c.TemplateData))
}

func (b *BranchAction) IsReady(c *engine.Context) bool {
	return engine.IsTruthy(b.ShouldRunCond.Render(c.TemplateData))
}

func (b *BranchAction) String() string {
	return fmt.Sprintf("branching %s", b.GetIf())
}

func (b *BranchAction) Validate() error {
	for _, v := range b.CompiledSteps {
		err := v.Validate()
		if err != nil {
			return errors.Join(fmt.Errorf("branch %s", b.ID), err)
		}
	}

	return nil
}

func (b *BranchAction) ToProto() *v1beta1.Action_Branch {
	return b.Action_Branch
}

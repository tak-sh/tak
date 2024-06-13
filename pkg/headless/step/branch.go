package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
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

var _ Action = &BranchAction{}
var _ grpcutils.ProtoWrapper[*v1beta1.Action_Branch] = &BranchAction{}
var _ PathNode = &BranchAction{}

type BranchAction struct {
	*v1beta1.Action_Branch
	ID            string
	CompiledSteps []*Step
	ShouldRunCond *engine.TemplateRenderer
}

func (b *BranchAction) IsReady(st *engine.TemplateData) bool {
	return engine.IsTruthy(b.ShouldRunCond.Render(st))
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

func (b *BranchAction) Act(ctx *engine.Context) error {
	if !engine.IsTruthy(b.ShouldRunCond.Render(ctx.TemplateData)) {
		return nil
	}

	for _, v := range b.CompiledSteps {
		err := ctx.Evaluator.Eval(ctx, v)
		if err != nil {
			return errors.Join(fmt.Errorf("step %s", v.GetId()), err)
		}
	}

	return nil
}

func (b *BranchAction) GetID() string {
	return b.ID
}

func (b *BranchAction) ToProto() *v1beta1.Action_Branch {
	return b.Action_Branch
}

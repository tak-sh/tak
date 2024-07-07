package actionutils

import (
	"github.com/tak-sh/tak/pkg/mocks/enginemocks"
	"github.com/tak-sh/tak/pkg/mocks/stepmocks"
)

type Action struct {
	*stepmocks.Action
	*enginemocks.PathNode
}

type BranchAction struct {
	*stepmocks.Action
	*enginemocks.PathNode
	*stepmocks.Branches
}

func NewBranchAction() *BranchAction {
	return &BranchAction{
		Branches: new(stepmocks.Branches),
		Action:   new(stepmocks.Action),
		PathNode: new(enginemocks.PathNode),
	}
}

func NewAction() *Action {
	return &Action{
		Action:   new(stepmocks.Action),
		PathNode: new(enginemocks.PathNode),
	}
}

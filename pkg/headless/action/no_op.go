package action

import "github.com/tak-sh/tak/pkg/headless"

var _ Action = &NoOpAction{}

type NoOpAction struct {
	ID string
}

func (n *NoOpAction) Validate() error {
	return nil
}

func (n *NoOpAction) String() string {
	return "none"
}

func (n *NoOpAction) Act(_ *headless.Context) error { return nil }

func (n *NoOpAction) GetID() string { return n.ID }

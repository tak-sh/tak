package account

import (
	"context"
)

var _ Account = &ChaseBank{}

type ChaseBank struct {
}

func (c *ChaseBank) Login(ctx context.Context, username, password string) error {
	return nil
}

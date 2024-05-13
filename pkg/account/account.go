package account

import "context"

type Account interface {
	Login(ctx context.Context, username, password string) error
}

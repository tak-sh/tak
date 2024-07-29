package provider

import (
	"context"
	"github.com/tak-sh/tak/generated/go/api/provider/v1beta1"
)

type Provider interface {
	ListAccounts(ctx context.Context) ([]*v1beta1.Account, error)
	Login(ctx context.Context) error
	DownloadTransactions(ctx context.Context, accountName string) error
}

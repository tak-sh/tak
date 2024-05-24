package account

import (
	"errors"
	"github.com/tak-sh/tak/generated/go/api/account/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/protoenc"
	"os"
	"path/filepath"
)

func LoadFile(fp string) (*v1beta1.Account, error) {
	_, err := os.Stat(fp)
	if err != nil {
		return nil, errors.Join(except.NewNotFound("failed to find account file %s", fp), err)
	}

	acct := new(v1beta1.Account)
	dir, name := filepath.Split(fp)
	err = protoenc.UnmarshalFile(acct, name, os.DirFS(dir))
	if err != nil {
		return nil, errors.Join(except.NewInvalid("%s is not a valid account file", fp), err)
	}

	return acct, nil
}

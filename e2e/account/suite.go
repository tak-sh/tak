package account

import (
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/account/v1beta1"
	"github.com/tak-sh/tak/pkg/account"
	"github.com/tak-sh/tak/pkg/utils/fileutils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

type Suite struct {
	suite.Suite
}

func (s *Suite) StartAccountTest(accountName, fileName string) (*TestRun, error) {
	accountsPath := fileutils.FindUpwardFrom("accounts", "", "")
	chasePath := filepath.Join(accountsPath, accountName)
	acct, err := account.LoadFile(filepath.Join(chasePath, fileName))
	if err != nil {
		return nil, err
	}

	out := &TestRun{
		Account: acct,
	}
	out.Server = httptest.NewServer(http.FileServerFS(os.DirFS(filepath.Join(chasePath, "html"))))

	return out, nil
}

type TestRun struct {
	Account *v1beta1.Account
	Server  *httptest.Server
}

func (t *TestRun) Close() {
	if t.Server != nil {
		t.Server.Close()
	}
}

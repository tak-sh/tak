package account

import (
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/provider/v1beta1"
	"github.com/tak-sh/tak/pkg/provider"
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
	acct, err := provider.LoadFile(filepath.Join(chasePath, fileName))
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
	Account *v1beta1.Provider
	Server  *httptest.Server
}

func (t *TestRun) Close() {
	if t.Server != nil {
		t.Server.Close()
	}
}

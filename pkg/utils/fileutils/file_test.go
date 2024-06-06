package fileutils

import (
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type FileTestSuite struct {
	suite.Suite
}

func (f *FileTestSuite) TestFindUpwardFrom() {
	type test struct {
		Name     string
		Start    string
		End      string
		Paths    []string
		Expected string
	}

	tests := map[string]test{
		"happy path": {
			Name:  "accounts",
			Start: filepath.Join("derp", "accounts", "a", "b"),
			Paths: []string{
				filepath.Join("derp", "accounts", "a", "b"),
				filepath.Join("derp", "accounts", "a", "b"),
			},
			Expected: filepath.Join("derp", "accounts"),
		},
		"not_found": {
			Name:  "account",
			Start: filepath.Join("derp", "accounts", "a", "b"),
			Paths: []string{
				filepath.Join("derp", "accounts", "a", "b"),
				filepath.Join("derp", "accounts", "a", "b"),
			},
		},
	}

	for descr, t := range tests {
		func() {
			basePath := strings.ReplaceAll(descr, " ", "_")
			base := filepath.Join(os.TempDir(), basePath)
			for _, p := range t.Paths {
				var tempDir string
				isFile := filepath.Ext(p) != ""
				if isFile {
					tempDir = filepath.Join(base, filepath.Dir(p))
				} else {
					tempDir = filepath.Join(base, p)
				}
				_ = os.MkdirAll(tempDir, os.ModePerm)
				if isFile {
					_ = os.WriteFile(p, []byte{}, os.ModePerm)
				}
			}
			defer os.RemoveAll(base)

			t.Start = absPath(base, t.Start)
			t.End = absPath(base, t.End)
			t.Expected = absPath(base, t.Expected)

			actual := FindUpwardFrom(t.Name, t.Start, t.End)
			f.Equal(t.Expected, actual, descr)
		}()
	}
}

func absPath(base, fp string) string {
	if !filepath.IsAbs(fp) && fp != "" {
		fp = filepath.Join(base, fp)
	}
	return fp
}

func TestFileTestSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}

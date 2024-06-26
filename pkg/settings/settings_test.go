package settings

import (
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/settings/v1beta1"
	"github.com/tak-sh/tak/pkg/protoenc"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"github.com/tak-sh/tak/pkg/utils/testutils"
	"os"
	"testing"
	"testing/fstest"
)

type SettingsTestSuite struct {
	suite.Suite
}

func (s *SettingsTestSuite) TestGet() {
	type test struct {
		EnvVars  map[string]string
		Initial  *v1beta1.Settings
		Expected *v1beta1.Settings
	}

	tests := map[string]test{
		"obeys env var precedence": {
			EnvVars: map[string]string{
				"TAK_CHROMEDATADIRECTORY": "FLERP",
			},
			Initial: &v1beta1.Settings{
				ChromeDataDirectory: ptr.Ptr("derp"),
			},
			Expected: &v1beta1.Settings{
				ChromeDataDirectory: ptr.Ptr("FLERP"),
			},
		},
		"obeys env var precedence with bad file": {
			EnvVars: map[string]string{
				"TAK_CHROMEDATADIRECTORY": "FLERP",
			},
			Expected: &v1beta1.Settings{
				ChromeDataDirectory: ptr.Ptr("FLERP"),
			},
		},
		"defaults properly": {
			Expected: &v1beta1.Settings{
				ChromeDataDirectory: ptr.Ptr("chrome_user_data"),
			},
		},
	}

	for desc, t := range tests {
		func() {
			for k, v := range t.EnvVars {
				_ = os.Setenv(k, v)
			}

			defer func() {
				for k := range t.EnvVars {
					_ = os.Unsetenv(k)
				}
			}()

			tFs := fstest.MapFS{}
			set := New("", "", "").(*settings)
			if t.Initial != nil {
				b, _ := protoenc.MarshalYAML(t.Initial)
				tFs[FileName] = &fstest.MapFile{Data: b}
			}

			set.SetFS = tFs
			s.Empty(testutils.DiffProto(t.Expected, set.Get()), desc)
		}()
	}
}

func TestSettingsTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsTestSuite))
}

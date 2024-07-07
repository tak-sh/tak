package step

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"testing"
)

type ActionTestSuite struct {
	suite.Suite
}

func (a *ActionTestSuite) TestAct() {
	type test struct {
		Given           *v1beta1.Action
		Id              string
		Post            func(desc string)
		ExpectedErr     string
		ExpectedCompErr string
		Ctx             *engine.Context
	}

	tests := map[string]test{
		"store works": func() test {
			c, _ := engine.NewContext(context.Background(), nil, nil, engine.ContextOpts{})

			return test{
				Given: &v1beta1.Action{
					Store: &v1beta1.Action_Store{
						KeyVals: []*v1beta1.KeyVal{
							{Key: "derp.flerp", Value: "1"},
						},
					},
				},
				Id: "derpflorp",
				Post: func(desc string) {
					o := c.TemplateData.GetStepVal("derpflorp.derp.flerp")
					a.Equal("1", o, desc)
				},
				Ctx: c,
			}
		}(),
	}

	for desc, t := range tests {
		given, err := New(t.Id, t.Given)
		if err != nil {
			if t.ExpectedCompErr != "" {
				a.EqualError(err, t.ExpectedCompErr, desc)
			} else {
				a.NoError(err)
			}
			continue
		}

		err = given.Act(t.Ctx)
		if err != nil {
			if t.ExpectedErr != "" {
				a.EqualError(err, t.ExpectedErr, desc)
			} else {
				a.NoError(err)
			}
			continue
		}

		if t.Post != nil {
			t.Post(desc)
		}
	}
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(ActionTestSuite))
}

package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"google.golang.org/protobuf/proto"
	"net/url"
	"time"
)

func NewNav(id string, a *v1beta1.Action_Nav) *Nav {
	out := &Nav{
		Action_Nav: a,
		ID:         id,
	}

	return out
}

var _ Action = &Nav{}

type Nav struct {
	*v1beta1.Action_Nav
	ID string
}

func (n *Nav) Message() proto.Message {
	return n.Action_Nav
}

func (n *Nav) GetId() string {
	return n.ID
}

func (n *Nav) Eval(c *engine.Context, to time.Duration) error {
	c, cancel := c.WithTimeout(to)
	defer cancel()
	return c.Browser.Navigate(c, n.GetAddr())
}

func (n *Nav) Cancel(err error) {
	//TODO implement me
	panic("implement me")
}

func (n *Nav) Validate() error {
	_, err := url.Parse(n.GetAddr())
	if err != nil {
		return errors.Join(except.NewInvalid("%s is not a valid addr", n.GetAddr()), err)
	}

	return nil
}

func (n *Nav) String() string {
	return fmt.Sprintf("navigating to %s", n.GetAddr())
}

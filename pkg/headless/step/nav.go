package step

import (
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"net/url"
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

func (n *Nav) Act(ctx *engine.Context) error {
	return chromedp.Navigate(n.GetAddr()).Do(ctx)
}

func (n *Nav) GetID() string {
	return n.ID
}

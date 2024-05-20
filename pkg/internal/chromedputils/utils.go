package chromedputils

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// SendKeysLazy is a copy of chromedp.SendKeys but with the ability to generate the string at runtime.
func SendKeysLazy(sel any, v func() string, opts ...chromedp.QueryOption) chromedp.QueryAction {
	return chromedp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		n := nodes[0]

		// grab type attribute from node
		typ, attrs := "", n.Attributes
		n.RLock()
		for i := 0; i < len(attrs); i += 2 {
			if attrs[i] == "type" {
				typ = attrs[i+1]
			}
		}
		n.RUnlock()

		// when working with input[type="file"], call dom.SetFileInputFiles
		if n.NodeName == "INPUT" && typ == "file" {
			return dom.SetFileInputFiles([]string{v()}).WithNodeID(n.NodeID).Do(ctx)
		}

		return chromedp.KeyEventNode(n, v()).Do(ctx)
	}, append(opts, chromedp.NodeVisible)...)
}

package cdp

import (
	"context"
	"errors"
	"fmt"
	"time"

	proto "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// Step executes the actions given and records how long it takes
func (cdp *Engine) Step(name string, actions ...chromedp.Action) {
	cdp.muStep.Lock()
	defer cdp.muStep.Unlock()
	cdp.stepInfo.start(name)
	defer cdp.stepInfo.end(name)
	if err := chromedp.Run(cdp.browser, actions...); err != nil {
		if errors.Is(err, context.Canceled) {
			err = fmt.Errorf("timeout %v", cdp.timeout)
		}
		cdp.ErrorScreenshot(err)
		panic(err)
	}
}

// IsPresent checks if something is present
func (cdp *Engine) IsPresent(sel interface{}, opts ...chromedp.QueryOption) bool {
	cdp.muStep.Lock()
	defer cdp.muStep.Unlock()
	p := make(chan bool)
	go func() {
		nodes := new([]*proto.Node)
		if err := chromedp.Run(cdp.browser, chromedp.Nodes(sel, nodes, opts...)); err != nil {
			if errors.Is(err, context.Canceled) {
				err = fmt.Errorf("timeout %v", cdp.timeout)
				p <- false
				return
			}
			cdp.ErrorScreenshot(err)
			p <- false
			return
		}

		p <- len(*nodes) > 0
	}()

	select {
	case res := <-p:
		return res
	case <-time.After(100 * time.Millisecond):
		return false
	}
}

// Body is used to check the content of the page
func (cdp *Engine) Body(checks ...szenario.CheckFunc) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var body string
		if err := chromedp.Run(cdp.browser, chromedp.Text(`/html/body`, &body)); err != nil {
			return err
		}
		var errTot error
		for i := 1; i < 5; i++ {
			errTot = nil
			for _, f := range checks {
				if err := f(&body); err != nil {
					errTot = err
				}
			}
			if errTot == nil {
				return nil
			}
			time.Sleep(time.Duration(i) * 100 * time.Millisecond)
		}
		return errTot
	})
}

// WaitForEver blocks until the timeout is reached
func (cdp *Engine) WaitForEver() {
	chromedp.Run(cdp.browser, chromedp.WaitReady(`#ThisWillNotBeFoundAndWeWaitForEver`, chromedp.ByID))
}

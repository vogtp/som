package cdp

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/core/log"
)

// IsHeadless indicates if the browser is headless (i.e. does not show on screen)
func (cdp *Engine) IsHeadless() bool {
	return cdp.headless
}

// Log returns the logger
func (cdp *Engine) Log() *slog.Logger {
	return cdp.log
}

func (cdp *Engine) createEngine() (cancel context.CancelFunc) {
	if cdp.browser != nil {
		if pc := chromedp.FromContext(cdp.browser); pc != nil && pc.Browser != nil {
			slog.Debug("Engine already Initialised")
			//	return pc.Cancel
		}
	}
	ctx := cdp.ctx //context.Background()
	//if !cdp.headless {
	ctx, _ = chromedp.NewExecAllocator(ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", cdp.headless),
			chromedp.Flag("incognito", true),
			chromedp.Flag("disable-background-networking", false),
		)...,
	)
	//}

	cdp.browser, cancel = chromedp.NewContext(
		ctx,
		chromedp.WithErrorf(func(msg string, v ...interface{}) { cdp.baseLogger.Error(fmt.Sprintf(msg, v)) }),
		chromedp.WithLogf(func(msg string, v ...interface{}) { cdp.baseLogger.Info(fmt.Sprintf(msg, v)) }),
		chromedp.WithDebugf(func(msg string, v ...interface{}) { cdp.baseLogger.Debug(fmt.Sprintf(msg, v)) }),
	)

	cdp.registerConsoleListener()

	if cdp.noClose {
		cancel = func() { slog.Info("No close requested keeping the window open!") }
	}
	return cancel
}

func (cdp *Engine) clearConsoleCounter() {
	cdp.consMsg["exception"] = 0
	cdp.consMsg["error"] = 0
	cdp.consMsg["warning"] = 0
	cdp.consMsg["assert"] = 0
}

func (cdp *Engine) registerConsoleListener() {
	cdp.clearConsoleCounter()
	chromedp.ListenTarget(cdp.browser, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			// log, debug, info, error, warning, dir, dirxml, table, trace, clear, startGroup, startGroupCollapsed, endGroup, assert, profile, profileEnd, count, timeEnd
			switch ev.Type {
			case "error":
				cdp.consMsg[string(ev.Type)] = cdp.consMsg[string(ev.Type)] + 1
			case "warning":
				cdp.consMsg[string(ev.Type)] = cdp.consMsg[string(ev.Type)] + 1
			case "assert":
				cdp.consMsg[string(ev.Type)] = cdp.consMsg[string(ev.Type)] + 1
			}

			// cdp.log.Debug("console.%s message:", ev.Type)
			// for _, arg := range ev.Args {
			// 	cdp.log.Debug("%s - %s", arg.Type, arg.Value)
			// }
		case *runtime.EventExceptionThrown:
			cdp.consMsg["exception"] = cdp.consMsg["exception"] + 1
			cdp.log.Debug("Console exception", log.Error, ev.ExceptionDetails.Text)

		}
	})
}

package cdp

import (
	"time"

	"github.com/vogtp/som/cmd/somctl/term"
	"github.com/vogtp/som/pkg/env"
)

func (cdp *Engine) endStepActions() {
	if cdp.stepBreakPoint != nil && env.IsGoRun() {
		cdp.stepBreakPoint <- cdp.stepInfo.name
		return
	}
	time.Sleep(cdp.stepDelay)
}

func (cdp *Engine) BreakWaitForUserInput() {
	if !env.IsGoRun() {
		return
	}
	term.Read("Press any key to continue...\n")
}

package cdp

import (
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/cmd/somctl/term"
)

func (cdp *Engine) endStepActions() {
	if cdp.stepBreakPoint != nil && hcl.IsGoRun() {
		cdp.stepBreakPoint <- cdp.stepInfo.name
		return
	}
	time.Sleep(cdp.stepDelay)
}

func (cdp *Engine) BreakWaitForUserInput() {
	if !hcl.IsGoRun() {
		return
	}
	term.Read("Press any key to continue...\n")
}

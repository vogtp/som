package cdp

import (
	"fmt"
	"time"
)

func (cdp *Engine) report(totalDuration time.Duration) {
	if !cdp.sendReport {
		cdp.hcl.Warn("Not sending reports probably killed by OS")
		return
	}
	cdp.mu.Lock()
	defer cdp.mu.Unlock()
	status := "OK"
	if cdp.evtMsg.Err() != nil {
		status = cdp.evtMsg.Err().Error()
	}
	cdp.evtMsg.SetCounter("step.total", totalDuration.Seconds())
	failedLogins := cdp.szenario.User().FailedLogins()
	cdp.evtMsg.SetCounter("logins.failed", float64(failedLogins))
	if failedLogins > 0 {
		pwAge := time.Since(cdp.szenario.User().PasswordCreated())
		cdp.evtMsg.SetCounter("logins.passwordage", float64(pwAge.Seconds()))
		cdp.evtMsg.SetStatus("logins.passwordage", fmt.Sprintf("%v", pwAge))
		cdp.hcl.Errorf("Failed logins: %v", failedLogins)
		cdp.hcl.Errorf("Password Age: %v", pwAge)
	}
	for k, v := range cdp.stepInfo.stepTimes {
		if v > 0 {
			cdp.evtMsg.SetCounter("step."+k, v)
		}
	}
	for k, v := range cdp.consMsg {
		if v > 0 {
			cdp.evtMsg.SetCounter("console."+k, float64(v))
		}
	}
	if err := cdp.bus.Szenario.Send(cdp.evtMsg); err != nil {
		cdp.hcl.Warnf("cannot send szenario message: %v", err)
	}
	cdp.hcl.Infof("Status %s: %v", cdp.szenario.Name(), status)
}

package cdp

import (
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
	pwAge := time.Since(cdp.szenario.User().PasswordCreated())
	cdp.evtMsg.SetCounter("logins.failed", float64(failedLogins))
	cdp.evtMsg.SetCounter("logins.passwordage", float64(pwAge.Seconds()))
	cdp.evtMsg.SetStatus("logins.passwordage", pwAge.String())
	cdp.hcl.Warn("Failed logins", "failed_login", failedLogins, "password_age", pwAge)

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
		cdp.hcl.Warn("cannot send szenario message", "error", err)
	}
	cdp.hcl.Info("Szenario status", status)
}

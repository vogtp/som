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
	cdp.bus.Szenario.Send(cdp.evtMsg)
	cdp.hcl.Infof("Status %s: %v", cdp.szenario.Name(), status)
}

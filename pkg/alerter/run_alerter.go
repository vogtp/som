package alerter

import (
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
)

// Run the alerter
func Run(name string, coreOpts ...core.Option) (func(), error) {
	c, close := core.New(name, coreOpts...)
	hcl := c.HCL()
	m, err := NewMailer()
	if err == nil {
		AddEngine(m)
	} else {
		hcl.Warnf("Cannot creater mailer engine: %v", err)
	}
	t, err := NewTeams()
	if err == nil {
		AddEngine(t)
	} else {
		hcl.Warnf("Cannot creater teams engine: %v", err)
	}

	parseConfig()
	c.Bus().Alert.Handle(handle)
	if teams, ok := t.(*Teams); ok {
		teams.checkDestinationWebhooks()
	}
	return close, nil
}

func parseConfig() {
	parseDestinationsCfg()
	parseRulesCfg()
}

func handle(a *msg.AlertMsg) {
	for _, r := range rules {
		// TODO check requirement
		for _, d := range r.Destinations {
			engines[d.Kind].Send(a, &d)
		}
	}
}

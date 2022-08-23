package alerter

import (
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/stater/alertmgr"
)

const (
	cfgAlertDestName         = "name"
	cfgAlertDestMailTo       = "to"
	cfgAlertDestTeamsWebhook = "webhook"
)

type Alerter struct {
	hcl     hcl.Logger
	dsts    map[string]*Destination
	engines map[string]Engine
	rules   []Rule
}

func New(c *core.Core) *Alerter {
	a := &Alerter{
		hcl:     c.HCL().Named("alerter"),
		dsts:    make(map[string]*Destination, 0),
		engines: make(map[string]Engine),
		rules:   make([]Rule, 0),
	}

	return a
}

func (a *Alerter) handle(msg *msg.AlertMsg) {
	for _, r := range a.rules {
		// TODO check requirement
		for _, d := range r.Destinations {
			a.engines[d.Kind].Send(msg, &d)
		}
	}
}

func (a *Alerter) parseConfig() {
	a.parseDestinationsCfg()
	a.parseRulesCfg()
}

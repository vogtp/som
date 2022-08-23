package alerter

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
)

const (
	cfgAlertDest             = "destinations"
	cfgAlertSubject          = "subject"
	cfgAlerMailFrom          = "from"
	cfgAlertDestName         = "name"
	cfgAlertDestMailTo       = "to"
	cfgAlertDestTeamsWebhook = "webhook"
)

func getCfgString(key string, r *Rule, d *Destination) string {
	v := viper.GetString(fmt.Sprintf("alert.%s", key))
	if s := d.Cfg.GetString(key); len(s) > 0 {
		v = s
	}
	if s := r.Cfg.GetString(key); len(s) > 0 {
		v = s
	}
	return v
}

// Alerter is the main alerter stuct
type Alerter struct {
	hcl     hcl.Logger
	c       *core.Core
	dsts    map[string]*Destination
	engines map[string]Engine
	rules   []Rule
}

// New creates an alerter
func New(c *core.Core) *Alerter {
	a := &Alerter{
		hcl:     c.HCL().Named("alerter"),
		c:       c,
		dsts:    make(map[string]*Destination, 0),
		engines: make(map[string]Engine),
		rules:   make([]Rule, 0),
	}
	return a
}

// Run the alerter
func (a *Alerter) Run() (ret error) {
	if err := a.initDests(); err != nil {
		ret = err
		a.hcl.Warnf("problems initialiseing alerter destinations: %v", err)
	}
	if err := a.initRules(); err != nil {
		ret = err
		a.hcl.Warnf("problems initialiseing alerter rules: %v", err)
	}
	if err := a.initEgninges(); err != nil {
		ret = err
		a.hcl.Warnf("problems initialiseing alerter engines: %v", err)
	}
	a.c.Bus().Alert.Handle(a.handle)
	return ret
}

func (a *Alerter) initEgninges() (ret error) {
	for _, e := range a.engines {
		if err := e.checkConfig(a); err != nil {
			a.hcl.Warnf("Engine %s has config errors: %v", e.Kind(), err)
			ret = err
		}
	}
	return ret
}

func (a *Alerter) handle(msg *msg.AlertMsg) {
	for _, r := range a.rules {
		// TODO check requirement
		for _, d := range r.Destinations {
			a.engines[d.Kind].Send(msg, &r, &d)
		}
	}
}

func (a *Alerter) parseConfig() {
	a.parseDestinationsCfg()
	a.parseRulesCfg()
}

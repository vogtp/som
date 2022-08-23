package alerter

import (
	"fmt"

	"github.com/spf13/cast"
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
	cfgAlertEnabled          = "enabled"
	cfgAlertRuleConditions   = "conditions"
	cfgInclude               = "include"
	cfgExclude               = "exclude"
)

func getCfgBool(key string, r *Rule, d *Destination) bool {
	c := getCfg(key, r, d)
	if c == nil {
		return true
	}
	return cast.ToBool(c)
}

func getCfgString(key string, r *Rule, d *Destination) string {
	return cast.ToString(getCfg(key, r, d))
}

func getCfg(key string, r *Rule, d *Destination) any {
	v := viper.Get(fmt.Sprintf("alert.%s", key))
	if d != nil {
		if s := d.Cfg.Get(key); s != nil {
			v = s
		}
	}
	if r != nil {
		if s := r.Cfg.Get(key); s != nil {
			v = s
		}
	}
	return v
}

// Alerter is the main alerter stuct
type Alerter struct {
	hcl        hcl.Logger
	c          *core.Core
	dsts       map[string]*Destination
	engines    map[string]Engine
	conditions map[string]Conditon
	rules      []Rule
}

// New creates an alerter
func New(c *core.Core) *Alerter {
	a := &Alerter{
		hcl:        c.HCL().Named("alerter"),
		c:          c,
		dsts:       make(map[string]*Destination),
		engines:    make(map[string]Engine),
		conditions: make(map[string]Conditon),
		rules:      make([]Rule, 0),
	}
	a.addDefaultComponents()
	return a
}

func (a *Alerter) addDefaultComponents() {
	a.AddConditon(StatusCond{})
	a.AddConditon(SzenarioCond{})
	if err := a.AddEngine(NewMailer()); err != nil {
		a.hcl.Warnf("Cannot create engine: %v", err)
	}
	if err := a.AddEngine(NewTeams()); err != nil {
		a.hcl.Warnf("Cannot create engine: %v", err)
	}
}

// Run the alerter
func (a *Alerter) Run() (ret error) {
	a.parseConfig()
	if err := a.initDests(); err != nil {
		ret = err
		a.hcl.Warnf("problems initialising alerter destinations: %v", err)
	}
	if err := a.initRules(); err != nil {
		ret = err
		a.hcl.Warnf("problems initialising alerter rules: %v", err)
	}
	if err := a.initEgninges(); err != nil {
		ret = err
		a.hcl.Warnf("problems initialising alerter engines: %v", err)
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
		if err := r.DoAlert(msg); err != nil {
			a.hcl.Infof("Not alerting %s: %v", msg.Name, err)
			continue
		}
		for _, d := range r.Destinations {
			if !getCfgBool(cfgAlertEnabled, &r, &d) {
				a.hcl.Warnf("not alerting %s alerting is disabled", msg.Name)
				continue
			}
			a.engines[d.Kind].Send(msg, &r, &d)
		}
	}
}

func (a *Alerter) parseConfig() {
	a.parseDestinationsCfg()
	a.parseRulesCfg()
}

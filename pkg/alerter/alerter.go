package alerter

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"golang.org/x/exp/slog"
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
		if s := d.cfg.Get(key); s != nil {
			v = s
		}
	}
	if r != nil {
		if s := r.cfg.Get(key); s != nil {
			v = s
		}
	}
	return v
}

// Alerter is the main alerter stuct
type Alerter struct {
	log        *slog.Logger
	c          *core.Core
	dsts       map[string]*Destination
	engines    map[string]Engine
	conditions map[string]Conditon
	rules      []Rule
}

// New creates an alerter
func New(c *core.Core) *Alerter {
	a := &Alerter{
		log:        c.HCL().With(log.Component, "alerter"),
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
	if err := a.AddConditon(StatusCond{}); err != nil {
		a.log.Warn("Cannot add status condition", log.Error, err)
	}
	if err := a.AddConditon(SzenarioCond{}); err != nil {
		a.log.Warn("Cannot add szenario condition", log.Error, err)
	}
	if err := a.AddEngine(NewMailer()); err != nil {
		a.log.Warn("Cannot create engine", log.Error, err)
	}
	if err := a.AddEngine(NewTeams()); err != nil {
		a.log.Warn("Cannot create engine", log.Error, err)
	}
}

// Run the alerter
func (a *Alerter) Run() (ret error) {
	a.parseConfig()
	if err := a.initDests(); err != nil {
		ret = err
		a.log.Warn("problems initialising alerter destinations", log.Error, err)
	}
	if err := a.initRules(); err != nil {
		ret = err
		a.log.Warn("problems initialising alerter rules", log.Error, err)
	}
	if err := a.initEgninges(); err != nil {
		ret = err
		a.log.Warn("problems initialising alerter engines", log.Error, err)
	}
	a.c.Bus().Alert.Handle(a.handle)
	return ret
}

func (a *Alerter) initEgninges() (ret error) {
	for _, e := range a.engines {
		if err := e.checkConfig(a); err != nil {
			a.log.Warn("Engine has config errors", "engine", e.Kind(), log.Error, err)
			ret = err
		}
	}
	return ret
}

func (a *Alerter) handle(msg *msg.AlertMsg) {
	for _, r := range a.rules {
		if err := r.DoAlert(msg); err != nil {
			a.log.Info("Not alerting", "alert", msg.Name, log.Error, err, "rule", r.name)
			continue
		}
		for _, d := range r.destinations {
			if !getCfgBool(cfgAlertEnabled, &r, &d) {
				a.log.Warn("alerting is disabled", "alert", msg.Name, "destination", d.name, "rule", r.name)
				continue
			}
			if err := a.engines[d.kind].Send(msg, &r, &d); err != nil {
				a.log.Error("Cannot send message", "engine", d.kind, "destination", d.name, log.Error, err, "rule", r.name)
			}
		}
	}
}

func (a *Alerter) parseConfig() {
	a.parseDestinationsCfg()
	a.parseRulesCfg()
}

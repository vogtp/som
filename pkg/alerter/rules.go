package alerter

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
)

// Rule a rule for alerting
type Rule struct {
	Name         string
	Destinations []Destination
	Conditions   []Conditon
	Cfg          *viper.Viper
}

// Check checks if the condtions a matched
func (r *Rule) Check(mgs *msg.AlertMsg) bool {
	for _, c := range r.Conditions {
		if !c.Check(mgs) {
			return false
		}
	}
	return true
}

// AddRule adds an alerting Rule
func (a *Alerter) AddRule(r *Rule) error {
	if r == nil {
		return errors.New("rule is nil")
	}
	if len(r.Name) < 1 {
		return errors.New("a rule must have an name")
	}
	if !getCfgBool(cfgAlertEnabled, r, nil) {
		return fmt.Errorf("%s is not enabled", r.Name)
	}
	a.rules = append(a.rules, *r)
	return nil
}

func (a *Alerter) initRules() (ret error) {
	validRules := make([]Rule, 0, len(a.rules))
	for _, r := range a.rules {
		if err := a.isValidRule(&r); err != nil {
			ret = fmt.Errorf("rule %s is not valid: %v", r.Name, err)
			a.hcl.Warn(ret.Error())
			ret = err
		}
		validRules = append(validRules, r)
	}
	a.rules = validRules
	if len(validRules) < 1 {
		ret = errors.New("no valid alerting rules")
		a.hcl.Error(ret.Error())
	}
	a.hcl.Warnf("Loaded %v alert Rules", len(a.rules))
	return ret
}

func (a *Alerter) isValidRule(r *Rule) error {
	dests := r.Cfg.GetStringSlice(cfgAlertDest)
	for _, d := range dests {
		dst, found := a.dsts[d]
		if !found {
			a.hcl.Warnf("No such destination %q ignroing", d)
			continue
		}
		r.Destinations = append(r.Destinations, *dst)
	}
	if len(r.Destinations) < 1 {
		return errors.New("a rule without destinations does not make sens")
	}
	a.hcl.Infof("Added rule %s", r.Name)
	return nil
}

func (a *Alerter) parseRulesCfg() {
	raw := viper.Get(cfg.AlertRules)
	slc, ok := raw.([]any)
	if !ok {
		a.hcl.Errorf("Cannot get rules: %v", raw)
		return
	}
	for i := range slc {
		cfg := viper.Sub(fmt.Sprintf("%s.%v", cfg.AlertRules, i))
		name := cfg.GetString(cfgAlertDestName)
		if len(name) < 1 {
			a.hcl.Warn("No destination name, skipping")
			continue
		}
		r := &Rule{
			Name: name,
			Cfg:  cfg,
		}
		if err := a.AddRule(r); err != nil {
			a.hcl.Warnf("Not adding rule %s: %v", name, err)
		}
	}
}

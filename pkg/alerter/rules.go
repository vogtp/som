package alerter

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
)

type condWrapper struct {
	cond Conditon
	cfg  *viper.Viper
}

// Rule a rule for alerting
type Rule struct {
	name         string
	destinations []Destination
	conditions   []condWrapper
	cfg          *viper.Viper
}

// DoAlert checks if the condtions a matched
func (r *Rule) DoAlert(mgs *msg.AlertMsg) error {
	for _, c := range r.conditions {
		if err := c.cond.DoAlert(mgs, c.cfg); err != nil {
			return fmt.Errorf("rule %q condtion %q: %v", r.name, c.cond.Kind(), err)
		}
	}
	return nil
}

// AddRule adds an alerting Rule
func (a *Alerter) AddRule(r *Rule) error {
	if r == nil {
		return errors.New("rule is nil")
	}
	if len(r.name) < 1 {
		return errors.New("a rule must have an name")
	}
	if !getCfgBool(cfgAlertEnabled, r, nil) {
		return fmt.Errorf("%s is not enabled", r.name)
	}
	a.rules = append(a.rules, *r)
	return nil
}

func (a *Alerter) initRules() (ret error) {
	validRules := make([]Rule, 0, len(a.rules))
	for _, r := range a.rules {
		if err := a.isValidRule(&r); err != nil {
			ret = fmt.Errorf("rule %s is not valid: %v", r.name, err)
			a.hcl.Warn(ret.Error())
			ret = err
			continue
		}
		a.parseConditions(&r)
		validRules = append(validRules, r)
	}
	a.rules = validRules
	if len(validRules) < 1 {
		ret = errors.New("no valid alerting rules")
		a.hcl.Error(ret.Error())
	}
	a.hcl.Warn("Loaded alert Rules", "count", len(a.rules))
	return ret
}

func (a *Alerter) isValidRule(r *Rule) error {
	dests := r.cfg.GetStringSlice(cfgAlertDest)
	for _, d := range dests {
		dst, found := a.dsts[d]
		if !found {
			a.hcl.Warn("No such destination: ignroing", "destination", d, "rule", r.name)
			continue
		}
		r.destinations = append(r.destinations, *dst)
	}
	if len(r.destinations) < 1 {
		return fmt.Errorf("a rule %s without destinations does not make sense", r.name)
	}
	a.hcl.Info("Added rule", "rule", r.name)
	return nil
}

func (a *Alerter) parseRulesCfg() {
	raw := viper.Get(cfg.AlertRules)
	slc, ok := raw.([]any)
	if !ok {
		a.hcl.Error("Cannot get rules", "raw", raw)
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
			name:       name,
			cfg:        cfg,
			conditions: make([]condWrapper, 0),
		}
		if err := a.AddRule(r); err != nil {
			a.hcl.Warn("Not adding rule", "rule", name, "error", err)
		}
	}
}

func (a *Alerter) parseConditions(r *Rule) {
	raw := r.cfg.Get(cfgAlertRuleConditions)
	slc, ok := raw.(map[string]any)
	if !ok {
		a.hcl.Error("Cannot get conditions of rule", "rule", r.name, "conditions", raw)
		return
	}
	for n := range slc {
		cond, ok := a.conditions[n]
		if !ok {
			a.hcl.Warn("rule: no such codition", "rule", r.name, "condition", n)
			continue
		}
		cfg := r.cfg.Sub(fmt.Sprintf("%s.%v", cfgAlertRuleConditions, n))
		if err := cond.CheckConfig(cfg); err != nil {
			a.hcl.Warn("Condition %q of rule %q contains errors: %v", "condition", cond.Kind(), "rule", r.name, "error", err)
		}
		r.conditions = append(r.conditions, condWrapper{
			cond: cond,
			cfg:  cfg,
		})
	}
}

package alerter

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

// Rule a rule for alerting
type Rule struct {
	Name         string
	Destinations []Destination
	Cfg          *viper.Viper
}

// AddRule adds an alerting Rule
func (a *Alerter) AddRule(r *Rule) error {
	if r == nil {
		return errors.New("rule is nil")
	}
	if len(r.Destinations) < 1 {
		return errors.New("a rule without destinations does not make sens")
	}
	a.hcl.Infof("Rule %s", r.Name)
	a.rules = append(a.rules, *r)
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
		dests := cfg.GetStringSlice("destinations")
		r := &Rule{
			Name:         name,
			Destinations: make([]Destination, 0, len(dests)),
			Cfg:          cfg,
		}
		for _, d := range dests {
			dst, found := a.dsts[d]
			if !found {
				a.hcl.Warnf("No such destination %q ignroing")
				continue
			}
			r.Destinations = append(r.Destinations, *dst)
		}
		if err := a.AddRule(r); err != nil {
			a.hcl.Warnf("Not adding rule %s: %v", name, err)
		}
	}
	a.hcl.Warnf("Loaded %v alert Rules", len(a.rules))
}

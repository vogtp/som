package alerter

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

var rules []Rule = make([]Rule, 0)

// Rule a rule for alerting
type Rule struct {
	Name         string
	Destinations []Destination
	Cfg          *viper.Viper
}

// AddRule adds an alerting Rule
func AddRule(r *Rule) error {
	if r == nil {
		return errors.New("rule is nil")
	}
	if len(r.Destinations) < 1 {
		return errors.New("a rule without destinations does not make sens")
	}
	hcl.Infof("Rule %s", r.Name)
	rules = append(rules, *r)
	return nil
}

func parseRulesCfg() {
	hcl := core.Get().HCL().Named("rules")
	raw := viper.Get(cfg.AlertRules)
	slc, ok := raw.([]any)
	if !ok {
		hcl.Errorf("Cannot get rules: %v", raw)
		return
	}
	for i := range slc {
		cfg := viper.Sub(fmt.Sprintf("%s.%v", cfg.AlertRules, i))
		name := cfg.GetString(cfgAlertDestName)
		if len(name) < 1 {
			hcl.Warn("No destination name, skipping")
			continue
		}
		dests := cfg.GetStringSlice("destinations")
		r := &Rule{
			Name:         name,
			Destinations: make([]Destination, 0, len(dests)),
			Cfg:          cfg,
		}
		for _, d := range dests {
			dst, found := destinations[d]
			if !found {
				hcl.Warnf("No such destination %q ignroing")
				continue
			}
			r.Destinations = append(r.Destinations, *dst)
		}
		if err := AddRule(r); err != nil {
			hcl.Warnf("Not adding rule %s: %v", name, err)
		}
	}
	hcl.Warnf("Loaded %v alert Rules", len(rules))
}

package alerter

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

// Destination is a endpoint for alerting (e.g. mail to a group)
type Destination struct {
	name string
	kind string
	cfg  *viper.Viper
}

// AddDestination adds a destination (mail group, chat root) to alerting
func (a *Alerter) AddDestination(d *Destination) error {
	if d == nil {
		return errors.New("destination is nil")
	}
	if len(d.name) < 1 {
		return errors.New("a destination must have an name")
	}
	if !getCfgBool(cfgAlertEnabled, nil, d) {
		return fmt.Errorf("%s is not enabled", d.name)
	}
	_, found := a.dsts[d.name]
	if found {
		return fmt.Errorf("doublicated destination name %s: no adding", d.name)
	}
	a.dsts[d.name] = d
	return nil
}

func (a *Alerter) initDests() (ret error) {
	validDst := make(map[string]*Destination, len(a.dsts))
	for k, d := range a.dsts {
		_, exists := a.engines[d.kind]
		if !exists {
			ret = fmt.Errorf("destination kind %s does not extist", d.kind)
			a.log.Warn(ret.Error())
			continue
		}
		validDst[k] = d
	}
	a.dsts = validDst
	if len(a.dsts) < 1 {
		ret = errors.New("no valid alerting destinations")
		a.log.Error("no valid alerting destinations", "error", ret.Error())
	}
	a.log.Warn("Loaded alert destinations", "count", len(a.dsts))
	return ret
}

func (a *Alerter) parseDestinationsCfg() {
	raw := viper.Get(cfg.AlertDestinations)
	slc, ok := raw.([]any)
	if !ok {
		a.log.Error("Cannot get destinations", "raw", raw)
		return
	}
	for i, l := range slc {
		m := l.(map[string]any)
		for k := range m {
			cfg := viper.Sub(fmt.Sprintf("%s.%v.%v", cfg.AlertDestinations, i, k))
			if cfg == nil {
				a.log.Warn("Destination is nil", "index", i, "engine", k)
				continue
			}
			name := cfg.GetString(cfgAlertDestName)
			if len(name) < 1 {
				a.log.Warn("No destination name, skipping")
				continue
			}
			a.log.Info("Alert destination", "engine", k, "name", name)
			d := &Destination{
				name: name,
				kind: k,
				cfg:  cfg,
			}

			if err := a.AddDestination(d); err != nil {
				a.log.Warn("Cannot add destination", "name", name, "error", err)
			}
		}
	}
}

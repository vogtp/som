package alerter

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

const (
	cfgAlertDestName         = "name"
	cfgAlertDestMailTo       = "to"
	cfgAlertDestTeamsWebhook = "webhook"
)

var destinations map[string]*Destination = make(map[string]*Destination, 0)

// Destination is a endpoint for alerting (e.g. mail to a group)
type Destination struct {
	Name string
	Kind string
	Cfg  *viper.Viper
}

// AddDestination adds a destination (mail group, chat root) to alerting
func AddDestination(d *Destination) error {
	if d == nil {
		return errors.New("destination is nil")
	}
	if len(d.Name) < 1 {
		return errors.New("a destination must have an name")
	}
	_, exists := engines[d.Kind]
	if !exists {
		return fmt.Errorf("destination kind %s does not extist", d.Kind)
	}
	_, found := destinations[d.Name]
	if found {
		return fmt.Errorf("doublicated destination name %s: no adding", d.Name)
	}

	destinations[d.Name] = d
	return nil
}

func parseDestinationsCfg() {
	hcl := core.Get().HCL().Named("destinations")
	raw := viper.Get(cfg.AlertDestinations)
	slc, ok := raw.([]any)
	if !ok {
		hcl.Errorf("Cannot get destinations: %v", raw)
		return
	}
	for i, l := range slc {
		m := l.(map[string]any)
		for k := range m {
			cfg := viper.Sub(fmt.Sprintf("%s.%v.%v", cfg.AlertDestinations, i, k))
			if cfg == nil {
				hcl.Warnf("Destination index %v is nil", i)
				continue
			}
			name := cfg.GetString(cfgAlertDestName)
			if len(name) < 1 {
				hcl.Warn("No destination name, skipping")
				continue
			}
			hcl.Infof("Alert destination %v: %q", k, name)
			d := &Destination{
				Name: name,
				Kind: k,
				Cfg:  cfg,
			}

			if err := AddDestination(d); err != nil {
				hcl.Warnf("Cannot add destination %s: %v", name, err)
			}
		}
	}
	hcl.Warnf("Loaded %v alert destinations", len(destinations))
}

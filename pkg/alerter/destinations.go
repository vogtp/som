package alerter

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

type destinationKind uint

const (
	kindUnknown destinationKind = iota
	kindMail
	kindTeams
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
	Kind destinationKind
	Cfg  *viper.Viper
}

func AddDestination(d *Destination) error {
	if d == nil {
		return errors.New("destination is nil")
	}
	if len(d.Name) < 1 {
		return errors.New("a destination must have an name")
	}

	_, found := destinations[d.Name]
	if found {
		return fmt.Errorf("doublicated destination name %s: no adding", d.Name)
	}

	destinations[d.Name] = d
	return nil
}

func parseDestinations() {
	hcl := core.Get().HCL().Named("destinations")
	raw := viper.Get(cfg.AlertDestinations)
	slc, ok := raw.([]any)
	if !ok {
		hcl.Errorf("Cannot get destinations: %v", raw)
		return
	}
	for i, l := range slc {
		m := l.(map[string]any)
		for k, _ := range m {
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
			kind, err := parseKind(k)
			if err != nil {
				hcl.Warnf("Cannot parse destination: %v", err)
				continue
			}
			hcl.Infof("Alert destination %v: %q", k, name)
			d := &Destination{
				Name: name,
				Kind: kind,
				Cfg:  cfg,
			}

			if err := AddDestination(d); err != nil {
				hcl.Warnf("Cannot add destination %s: %v", name, err)
			}
		}
	}
	hcl.Warnf("Loaded %v alert destinations", len(destinations))
}

func parseKind(k string) (destinationKind, error) {
	switch k {
	case "mail":
		return kindMail, nil
	case "teams":
		return kindTeams, nil
	default:
		return kindUnknown, fmt.Errorf("Unknown destination kind: %s", k)
	}
}

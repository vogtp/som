package alerter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

// Conditon is a check if a alerting rule is triggered
type Conditon interface {
	Kind() string
	DoAlert(*msg.AlertMsg, *viper.Viper) error
	CheckConfig(*viper.Viper) error
}

// AddConditon add an alerting condition
func (a *Alerter) AddConditon(c Conditon) error {
	if c == nil {
		return errors.New("condition must not be nil")
	}
	if len(c.Kind()) < 1 {
		return errors.New("condition must have a kind")
	}
	a.conditions[c.Kind()] = c
	return nil
}

// StatusCond is a filter condtion based on the status
type StatusCond struct{}

// Kind is the name
func (StatusCond) Kind() string { return "status" }

// CheckConfig checks if the config is valid
func (StatusCond) CheckConfig(cfg *viper.Viper) error {
	lvl := status.Unknown.FromString(cfg.GetString("level"))
	if lvl == status.Unknown {
		return fmt.Errorf("reqested unknown status level: %s", cfg.GetString("level"))
	}
	return nil
}

// DoAlert checks if an alert should be send and returns an error if not
func (StatusCond) DoAlert(msg *msg.AlertMsg, cfg *viper.Viper) error {
	lvl := status.Unknown.FromString(cfg.GetString("level"))
	if lvl == status.Unknown {
		return fmt.Errorf("not alerting, reqested unknown status level: %s", cfg.GetString("level"))
	}
	if status.Unknown.FromString(msg.Level) < lvl {
		return fmt.Errorf("%s < %s", status.Unknown.FromString(msg.Level), lvl)
	}
	return nil
}

// SzenarioCond is a filter condtion based on the szenario name
type SzenarioCond struct{}

// Kind is the name
func (SzenarioCond) Kind() string { return "szenario" }

// CheckConfig checks if the config is valid
func (SzenarioCond) CheckConfig(cfg *viper.Viper) error {
	var err error
	szCfg := core.Get().SzenaioConfig()
	if szCfg == nil {
		return nil
	}
	szList := cfg.GetStringSlice(cfgInclude)
	szList = append(szList, cfg.GetStringSlice(cfgExclude)...)
	for _, n := range szList {
		if sz := szCfg.ByName(n); sz == nil {
			if err == nil {
				err = errors.New(n)
			} else {
				err = fmt.Errorf("%w, %q", err, n)
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("no such szenario %w", err)
	}
	return err
}

// DoAlert checks if an alert should be send and returns an error if not
func (SzenarioCond) DoAlert(msg *msg.AlertMsg, cfg *viper.Viper) error {
	name := strings.ToLower(msg.Name)
	include := cfg.GetStringSlice(cfgInclude)
	var err error
	if len(include) > 0 {
		err = fmt.Errorf("szenario %s is not included", msg.Name)
	}
	for _, n := range include {
		if strings.ToLower(n) == name {
			return nil // is included
		}
	}
	if err != nil {
		return err
	}
	for _, n := range cfg.GetStringSlice(cfgExclude) {
		if strings.ToLower(n) == name {
			return fmt.Errorf("szenario %s is excluded", msg.Name)
		}
	}
	return nil
}

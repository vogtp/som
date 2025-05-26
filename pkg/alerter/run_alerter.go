package alerter

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
)

// Run the alerter
func Run(name string, coreOpts ...core.Option) (func(), error) {
	pflag.Bool(cfg.AlertEnabled, true, "Disable alerting")
	cfg.Parse()
	c, close := core.New(name, coreOpts...)
	if !viper.GetBool(cfg.AlertEnabled) {
		c.Log().Warn("Alerting is disabled!")
		return close, nil
	}
	a := New(c)

	if err := a.Run(); err != nil {
		a.log.Warn("Problems starting the alerter", log.Error, err)
	}

	return close, nil
}

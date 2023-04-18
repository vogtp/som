package alerter

import (
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

// Run the alerter
func Run(name string, coreOpts ...core.Option) (func(), error) {
	c, close := core.New(name, coreOpts...)
	if !viper.GetBool(cfg.AlertEnabled) {
		c.HCL().Warn("Alerting is disabled!")
		return close, nil
	}
	a := New(c)

	if err := a.Run(); err != nil {
		a.log.Warn("Problems starting the alerter", "error", err)
	}

	return close, nil
}

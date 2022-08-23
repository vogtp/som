package alerter

import (
	"github.com/vogtp/som/pkg/core"
)

// Run the alerter
func Run(name string, coreOpts ...core.Option) (func(), error) {
	c, close := core.New(name, coreOpts...)
	a := New(c)

	if err := a.AddEngine(NewMailer()); err != nil {
		a.hcl.Warnf("Cannot create engine: %v", err)
	}
	if err := a.AddEngine(NewTeams()); err != nil {
		a.hcl.Warnf("Cannot create engine: %v", err)
	}

	a.parseConfig()
	if err := a.Run(); err != nil {
		a.hcl.Warnf("Error running the alerter: %v", err)
	}

	return close, nil
}

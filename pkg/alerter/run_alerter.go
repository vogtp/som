package alerter

import (
	"github.com/vogtp/som/pkg/core"
)

// Run the alerter
func Run(name string, coreOpts ...core.Option) (func(), error) {
	c, close := core.New(name, coreOpts...)
	a := New(c)

	m, err := NewMailer()
	if err == nil {
		a.AddEngine(m)
	} else {
		a.hcl.Warnf("Cannot creater mailer engine: %v", err)
	}
	t, err := NewTeams()
	if err == nil {
		a.AddEngine(t)
	} else {
		a.hcl.Warnf("Cannot creater teams engine: %v", err)
	}

	a.parseConfig()
	if teams, ok := t.(*Teams); ok {
		teams.checkConfig(a)
	}
	c.Bus().Alert.Handle(a.handle)
	return close, nil
}

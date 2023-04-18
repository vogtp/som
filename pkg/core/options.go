package core

import (
	"github.com/hashicorp/go-hclog"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// Option configures the core
type Option func(*Core)

// WebPort sets the port of the webserver
func WebPort(p int) Option {
	return func(c *Core) {
		c.web.port = p
	}
}

// BasePath sets the root of the webserver path
func BasePath(s string) Option {
	return func(c *Core) {
		c.web.basepath = s
	}
}

// BusLogger enables and sets the loggin of the bus bus
func BusLogger(level hclog.Level) Option {
	return func(c *Core) {
		c.bus.busLogLevel = level
	}
}

// Szenario sets the szenario config
func Szenario(szCfg *szenario.Config) Option {
	return func(c *Core) {
		c.szCfg = szCfg
	}
}

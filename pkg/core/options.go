package core

import (
	"github.com/hashicorp/go-hclog"
	"github.com/vogtp/go-hcl"
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

// HCL sets the logger
func HCL(hcl hcl.Logger) Option {
	return func(c *Core) {
		c.hcl = hcl
	}
}

// BusLogger enables and sets the loggin of the grav bus
func BusLogger(level hclog.Level) Option {
	return func(c *Core) {
		c.bus.gravLogLevel = level
	}
}

// Szenario sets the szenario config
func Szenario(szCfg *szenario.Config) Option {
	return func(c *Core) {
		c.szCfg = szCfg
	}
}

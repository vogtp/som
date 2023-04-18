package core

import (
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

var (
	muCreateCore sync.Mutex
	c            *Core
)

// Core is the central structure
type Core struct {
	hcl   hcl.Logger
	szCfg *szenario.Config
	name  string

	bus *Bus
	web *WebServer
}

// New creates a new Cores and its cleanup function
func New(name string, opts ...Option) (*Core, func()) {
	muCreateCore.Lock()
	defer muCreateCore.Unlock()
	cfg.Parse()
	newCore := false
	if c == nil {
		newCore = true
		c = &Core{
			hcl:  hcl.New(hcl.WithName(name), cfg.HclOptions()),
			name: name,
			web: &WebServer{
				port:     viper.GetInt(cfg.WebPort),
				basepath: viper.GetString(cfg.WebURLBasePath),
			},
			bus: &Bus{
				busLogLevel: hclog.LevelFromString(viper.GetString(cfg.BusLogLevel)),
			},
		}
	} else if c.name != name {
		c.hcl.Error("Cannot have two cores of different names", "name", c.name, "new_name", name)
	}
	for _, o := range opts {
		o(c)
	}
	if newCore {
		c.hcl.Warn("SOM starting...", "version", som.Version)
		c.web.init(c)
		c.bus.init(c)
		c.web.Start()
	}

	waitDuration := viper.GetDuration(cfg.CoreStartdelay)
	c.hcl.Info("Waiting for the core to get started up", "duration", waitDuration)
	<-time.After(waitDuration)
	return c, c.cleanup
}

// Get returns the core instance or panics if not Initialised with New
func Get() *Core {
	if c == nil {
		panic("Core must be Initialised with New first")
	}
	return c
}

// Bus returns the bus or panics if Core not Initialised with New
func (c *Core) Bus() *Bus {
	return c.bus
}

// HCL returns the logger or panics if Core not Initialised with New
func (c *Core) HCL() hcl.Logger {
	return c.hcl
}

// WebServer returns the webserver
func (c *Core) WebServer() *WebServer {
	return c.web
}

// SzenaioConfig returns the szenario config
func (c *Core) SzenaioConfig() *szenario.Config {
	if c.szCfg == nil {
		return szenario.NoConfig
	}
	return c.szCfg
}

func (c *Core) cleanup() {
	c.bus.cleanup()
	c.web.Stop()
}

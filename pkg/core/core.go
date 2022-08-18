package core

import (
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
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
				gravLogLevel: hclog.LevelFromString(viper.GetString(cfg.BusLogLevel)),
			},
		}
	} else if c.name != name {
		c.hcl.Errorf("Cannot have two cores of different names: %s != %s", c.name, name)
	}
	for _, o := range opts {
		o(c)
	}
	if newCore {
		c.hcl.Errorf("SOM %s starting...", cfg.Version)
		c.web.init(c)
		c.bus.init(c)
		c.web.Start()
	}

	waitDuration := viper.GetDuration(cfg.CoreStartdelay)
	c.hcl.Infof("Waiting %v for the core to get started up", waitDuration)
	<-time.After(waitDuration)
	hcl.Debugf("Waited %v hopefully the bus is up and running", waitDuration)
	return c, c.cleanup
}

// Get returns the core instance or panics if not initalised with New
func Get() *Core {
	if c == nil {
		panic("Core must be initalised with New first")
	}
	return c
}

// Bus returns the bus or panics if Core not initalised with New
func (c *Core) Bus() *Bus {
	return c.bus
}

// HCL returns the logger or panics if Core not initalised with New
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

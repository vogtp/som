package core

import (
	"fmt"
	"sync"
	"time"

	"log/slog"

	"github.com/spf13/viper"
	"github.com/vogtp/som"
	"github.com/vogtp/som/pkg/core/bus"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

var (
	muCreateCore sync.Mutex
	c            *Core
)

// Core is the central structure
type Core struct {
	log   *slog.Logger
	szCfg *szenario.Config
	name  string

	bus *bus.Manager

	web *WebServer
}

// New creates a new Cores and its cleanup function
func New(name string, opts ...Option) (*Core, func(), error) {
	muCreateCore.Lock()
	defer muCreateCore.Unlock()
	newCore := false
	if c == nil {
		newCore = true
		c = &Core{
			log:  log.New("som"),
			name: name,
		}
	} else if c.name != name {
		c.log.Error("Cannot have two cores of different names", "name", c.name, "new_name", name)
	}
	for _, o := range opts {
		o(c)
	}
	if newCore {
		cfg.Parse()
		slog.SetDefault(c.log)
		c.log.Warn("SOM starting...", "version", som.Version)
	}

	waitDuration := viper.GetDuration(cfg.CoreStartdelay)
	c.log.Info("Waiting for the core to get started up", "duration", waitDuration)
	<-time.After(waitDuration)
	return c, c.cleanup, nil
}

// Get returns the core instance or panics if not Initialised with New
func Get() *Core {
	if c == nil {
		panic("Core must be Initialised with New first")
	}
	return c
}

// Bus returns the bus or panics if Core not Initialised with New
func (c *Core) Bus() *bus.Manager {
	if c.bus == nil {
		b, err := bus.New(c.log)
		if err != nil {
			panic(fmt.Errorf("initialising bus: %w", err))
		}
		c.bus = b
	}
	return c.bus
}

// Log returns the logger or panics if Core not Initialised with New
func (c *Core) Log() *slog.Logger {
	return c.log
}

// WebServer returns the webserver
func (c *Core) WebServer() *WebServer {
	if c.web != nil {
		return c.web
	}
	port := viper.GetInt(cfg.WebPort)
	basePath := viper.GetString(cfg.WebURLBasePath)
	c.web = &WebServer{
		log:      c.log.With(log.Component, "web", "port", port, "basePath", basePath),
		port:     port,
		basepath: basePath,
	}

	c.web.init(c)
	c.web.Start()
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
	if c.web != nil {
		c.web.Stop()
	}
	if c.bus != nil {
		c.bus.Close()
	}
}

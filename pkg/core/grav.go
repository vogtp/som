package core

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/suborbital/e2core/bus/bus"
	"github.com/suborbital/e2core/bus/discovery/local"
	"github.com/suborbital/e2core/bus/transport/websocket"
	"github.com/suborbital/vektor/vlog"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
)

func (e *Bus) initGrav(web *WebServer) {

	wsPath := viper.GetString(cfg.BusWsPath)
	e.hcl.Debugf("Bus is using endpoint: %s", wsPath)
	gwss := websocket.New()
	opts := []bus.OptionsModifier{
		bus.UseEndpoint(fmt.Sprintf("%v", web.port), wsPath),
		bus.UseMeshTransport(gwss),
		// bus.UseBelongsTo(*belong)
	}

	if e.busLogLevel == hclog.Off {
		opts = append(opts, bus.UseLogger(vlog.Default(vlog.Level("null"))))
	} else {
		hcl := e.hcl.Named("bus")
		hcl.SetLevel(e.busLogLevel)

		opts = append(opts, bus.UseLogger(hcl.Vlog()))
	}

	if !hcl.IsGoTest() {
		e.hcl.Info("Starting local discovery")
		locald := local.New()
		opts = append(opts, bus.UseDiscovery(locald))
	}

	e.bus = bus.New(opts...)

	if len(web.basepath) > 0 && web.basepath != "/" {
		// no baseurl here it is used internally
		web.mux.HandleFunc(wsPath, gwss.HTTPHandlerFunc())
	}
	web.HandleFunc(wsPath, gwss.HTTPHandlerFunc())

	e.hcl.Infof("Endpoints: %v", viper.GetStringSlice(cfg.BusEndpoints))
	for i, ep := range viper.GetStringSlice(cfg.BusEndpoints) {
		if err := e.bus.ConnectEndpoint(ep); err != nil {
			e.hcl.Warnf("Error connecting to endpoint %s: %v", ep, err)
			continue
		}
		e.hcl.Infof("Connected to %v endpoint %s", i, ep)
	}
}

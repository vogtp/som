package core

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/suborbital/grav/discovery/local"
	"github.com/suborbital/grav/grav"
	"github.com/suborbital/grav/transport/websocket"
	"github.com/suborbital/vektor/vlog"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
)

func (e *Bus) initGrav(web *WebServer) {

	wsPath := viper.GetString(cfg.BusWsPath)
	e.hcl.Debugf("Bus is using endpoint: %s", wsPath)
	gwss := websocket.New()
	opts := []grav.OptionsModifier{
		grav.UseEndpoint(fmt.Sprintf("%v", web.port), wsPath),
		grav.UseMeshTransport(gwss),
		// grav.UseBelongsTo(*belong)
	}

	if e.busLogLevel == hclog.Off {
		opts = append(opts, grav.UseLogger(vlog.Default(vlog.Level("null"))))
	} else {
		hcl := e.hcl.Named("grav")
		hcl.SetLevel(e.busLogLevel)

		opts = append(opts, grav.UseLogger(hcl.Vlog()))
	}

	if !hcl.IsGoTest() {
		e.hcl.Info("Starting local discovery")
		locald := local.New()
		opts = append(opts, grav.UseDiscovery(locald))
	}

	e.bus = grav.New(opts...)

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

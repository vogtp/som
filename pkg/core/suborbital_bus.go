package core

import (
	"fmt"

	"log/slog"

	"github.com/spf13/viper"
	"github.com/suborbital/grav/discovery/local"
	"github.com/suborbital/grav/grav"
	"github.com/suborbital/grav/transport/websocket"
	"github.com/suborbital/vektor/vlog"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/env"
)

func (e *Bus) initGrav(web *WebServer) {

	wsPath := viper.GetString(cfg.BusWsPath)
	e.log.Debug("Bus is using endpoint", "endpoint", wsPath)
	gwss := websocket.New()
	opts := []grav.OptionsModifier{
		grav.UseEndpoint(fmt.Sprintf("%v", web.port), wsPath),
		grav.UseMeshTransport(gwss),
		// grav.UseBelongsTo(*belong)
	}

	if e.busLogLevel > slog.LevelError {
		opts = append(opts, grav.UseLogger(vlog.Default(vlog.Level("null"))))
	} else {
		slog := log.Create("som.grav", e.busLogLevel)
		vlog := log.VlogCompat(slog)
		opts = append(opts, grav.UseLogger(vlog))
	}

	if !env.IsGoTest() {
		e.log.Info("Starting local discovery")
		locald := local.New()
		opts = append(opts, grav.UseDiscovery(locald))
	}

	e.bus = grav.New(opts...)

	if len(web.basepath) > 0 && web.basepath != "/" {
		// no baseurl here it is used internally
		web.mux.HandleFunc(wsPath, gwss.HTTPHandlerFunc())
	}
	web.HandleFunc(wsPath, gwss.HTTPHandlerFunc())

	e.log.Info("Init grav", "endpoints", viper.GetStringSlice(cfg.BusEndpoints))
	for _, ep := range viper.GetStringSlice(cfg.BusEndpoints) {
		if err := e.bus.ConnectEndpoint(ep); err != nil {
			e.log.Warn("Error connecting to endpoint", "endpoint", ep, log.Error, err)
			continue
		}
		e.log.Info("Connected to peer", "endpoint", ep)
	}
}

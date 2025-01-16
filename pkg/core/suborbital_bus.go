package core

import (
	"fmt"

	"log/slog"

	"github.com/spf13/viper"
	"github.com/suborbital/e2core/foundation/bus/bus"
	"github.com/suborbital/e2core/foundation/bus/discovery/local"
	"github.com/suborbital/e2core/foundation/bus/transport/websocket"
	"github.com/suborbital/vektor/vlog"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/env"
)

func (e *Bus) initGrav(web *WebServer) {

	wsPath := viper.GetString(cfg.BusWsPath)
	e.log.Debug("Bus is using endpoint", "endpoint", wsPath)
	gwss := websocket.New()
	opts := []bus.OptionsModifier{
		bus.UseEndpoint(fmt.Sprintf("%v", web.port), wsPath),
		bus.UseMeshTransport(gwss),
		// bus.UseBelongsTo(*belong)
	}

	if e.busLogLevel > slog.LevelError {
		opts = append(opts, bus.UseLogger(vlog.Default(vlog.Level("null"))))
	} else {
		slog := log.Create("som.bus", e.busLogLevel)
		vlog := log.VlogCompat(slog)
		opts = append(opts, bus.UseLogger(vlog))
	}

	if !env.IsGoTest() {
		e.log.Info("Starting local discovery")
		locald := local.New()
		opts = append(opts, bus.UseDiscovery(locald))
	}

	e.bus = bus.New(opts...)

	if len(web.basepath) > 0 && web.basepath != "/" {
		// no baseurl here it is used internally
		web.mux.HandleFunc(wsPath, gwss.HTTPHandlerFunc())
	}
	web.HandleFunc(wsPath, gwss.HTTPHandlerFunc())

	e.log.Info("Init bus", "endpoints", viper.GetStringSlice(cfg.BusEndpoints))
	for _, ep := range viper.GetStringSlice(cfg.BusEndpoints) {
		if err := e.bus.ConnectEndpoint(ep); err != nil {
			e.log.Warn("Error connecting to endpoint", "endpoint", ep, log.Error, err)
			continue
		}
		e.log.Info("Connected to peer", "endpoint", ep)
	}
}

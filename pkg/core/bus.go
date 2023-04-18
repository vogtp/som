package core

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/suborbital/grav/grav"
	"github.com/vogtp/go-hcl"
	mesh "github.com/vogtp/go-mesh"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/msgtype"
	"golang.org/x/exp/slog"
)

// Bus is the event bus
type Bus struct {
	log  *slog.Logger
	bus  *grav.Grav
	mesh *mesh.Mgr

	endpointURL string

	busLogLevel slog.Level

	Szenario *eventHandler[msg.SzenarioEvtMsg]
	Alert    *eventHandler[msg.AlertMsg]
	Incident *eventHandler[msg.IncidentMsg]
}

// newBus creates a new eventbus
func (e *Bus) init(c *Core) {
	e.log = c.log.With(log.Component, "bus")
	e.initGrav(c.web)
	e.Szenario = newHandler[msg.SzenarioEvtMsg](
		e.log,
		e.bus,
		msgtype.Event,
	)
	e.Alert = newHandler[msg.AlertMsg](
		e.log,
		e.bus,
		msgtype.Alert,
	)
	e.Incident = newHandler[msg.IncidentMsg](
		e.log,
		e.bus,
		msgtype.Incident,
	)
	e.endpointURL = fmt.Sprintf("%s%s", c.web.url, viper.GetString(cfg.BusWsPath))
	e.log.Info("Bus started", "endpoint", e.endpointURL)
	// FIXME remove hcl
	e.mesh = mesh.New(e.bus, &mesh.NodeConfig{
		Name:     c.name,
		Endpoint: e.endpointURL,
	}, mesh.Hcl(hcl.New()), mesh.ConnectPeers(true), mesh.BroadcastIntervall(5*time.Minute), mesh.Purge(5*time.Minute))
	c.web.HandleFunc("/bus", e.mesh.HandlerInfo)
}

// GetLogger returns the eventbus logger
func (e Bus) GetLogger() *slog.Logger {
	return e.log
}

func (e *Bus) cleanup() {
	e.Szenario.cleanup()
	e.Alert.cleanup()
	e.mesh.Stop()
	if err := e.bus.Withdraw(); err != nil {
		e.log.Warn("cannot withdraw from bus", log.Error, err)
	}
	if err := e.bus.Stop(); err != nil {
		e.log.Warn("cannot stop bus", log.Error, err)
	}
}

// WaitMsgProcessed waits until the managed cannels have their messages sent
func (e *Bus) WaitMsgProcessed() {
	e.Szenario.WaitMsgProcessed()
	e.Alert.WaitMsgProcessed()
}

// Connect to grav and return a pod
func (e *Bus) Connect() *grav.Pod {
	return e.bus.Connect()
}

// EndpointURL returns the URL the bus endpoint listens on
func (e Bus) EndpointURL() string {
	return e.endpointURL
}

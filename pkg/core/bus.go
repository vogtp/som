package core

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/suborbital/grav/grav"
	mesh "github.com/vogtp/go-grav-mesh"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/msgtype"
)

// Bus is the event bus
type Bus struct {
	hcl  hcl.Logger
	bus  *grav.Grav
	mesh *mesh.Mgr

	endpointURL string

	gravLogLevel hclog.Level

	Szenario *eventHandler[msg.SzenarioEvtMsg]
	Alert    *eventHandler[msg.AlertMsg]
	Incident *eventHandler[msg.IncidentMsg]
}

// newBus creates a new eventbus
func (e *Bus) init(c *Core) {
	e.hcl = c.hcl.Named("bus")
	e.initGrav(c.web)
	e.Szenario = newHandler[msg.SzenarioEvtMsg](
		e.hcl,
		e.bus,
		msgtype.Event,
	)
	e.Alert = newHandler[msg.AlertMsg](
		e.hcl,
		e.bus,
		msgtype.Alert,
	)
	e.Incident = newHandler[msg.IncidentMsg](
		e.hcl,
		e.bus,
		msgtype.Incident,
	)
	e.endpointURL = fmt.Sprintf("%s%s", c.web.url, viper.GetString(cfg.BusWsPath))
	e.hcl.Infof("Bus started: %s", e.endpointURL)
	e.mesh = mesh.New(e.bus, &mesh.NodeConfig{
		Name:     c.name,
		Endpoint: e.endpointURL,
	}, mesh.Hcl(e.hcl), mesh.ConnectPeers(true), mesh.BroadcastIntervall(5*time.Minute), mesh.Purge(5*time.Minute))
	c.web.HandleFunc("/bus", e.mesh.HandlerInfo)
}

// GetLogger returns the eventbus logger
func (e Bus) GetLogger() hcl.Logger {
	return e.hcl
}

func (e *Bus) cleanup() {
	e.Szenario.cleanup()
	e.Alert.cleanup()
	e.mesh.Stop()
	e.bus.Withdraw()
	e.bus.Stop()
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

//EndpointURL returns the URL the bus endpoint listens on
func (e Bus) EndpointURL() string {
	return e.endpointURL
}

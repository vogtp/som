package alertmgr

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"log/slog"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

const (
	// KeyTopology is the key for the topology status of alerts
	KeyTopology = "topology"
)

// Stater defines the state of an entity
type Stater interface {
	json.Marshaler
	json.Unmarshaler
	GetAlert(e msg.SzenarioEvtMsg) *msg.AlertMsg
}

// AlertMgr handles alerts
type AlertMgr struct {
	log            *slog.Logger
	bus            *core.Bus
	alertIntervall time.Duration
	alertLevel     status.Level
	mu             sync.Mutex
	basicStates    map[string]*basicState
	status         status.Status
	reopenTime     time.Duration
}

// Option configures AlertMgr
type Option func(*AlertMgr)

// AlertIntervall sets the intervall in which an alert is repeated
func AlertIntervall(i time.Duration) Option {
	return func(am *AlertMgr) {
		am.alertIntervall = i
	}
}

// New registers a alert manager on the event bus
func New(options ...Option) {
	bus := core.Get().Bus()
	am := AlertMgr{
		log:            bus.GetLogger().With(log.Component, "alertMgr"),
		bus:            bus,
		basicStates:    make(map[string]*basicState),
		status:         status.New(),
		alertIntervall: viper.GetDuration(cfg.AlertIntervall),
		alertLevel:     status.Unknown.FromString(viper.GetString(cfg.AlertLevel)),
		reopenTime:     viper.GetDuration(cfg.AlertIncidentCorrelationReopenTime),
	}
	if am.alertLevel == status.Unknown {
		panic(fmt.Sprintf("Unknown %s: %s", cfg.AlertLevel, viper.GetString(cfg.AlertLevel)))
	}
	if err := am.load(); err != nil {
		am.log.Error("Cannot read state", log.Error, err)
	}
	am.Configure(options...)
	bus.Szenario.Handle(am.handle)
	am.log = am.log.With("min_alert_level", am.alertLevel, "realert_interval", am.alertIntervall)
	am.log.Info("AlertMgr started")
}

// Configure the AlertMgr
func (am *AlertMgr) Configure(options ...Option) {
	for _, o := range options {
		o(am)
	}
}

func (am *AlertMgr) handle(e *msg.SzenarioEvtMsg) {
	defer func() {
		go func() {
			if err := am.save(); err != nil {
				am.log.Warn("Cannot save alertmgr", log.Error, err)
			}
		}()
	}()
	am.log.Debug("Got event", log.Szenario, e.Name, "message", e.Err())
	if a := am.checkEvent(e); a != nil {
		am.log.Warn("Generating alert for %v: %v (%v, %v)", e.Name, e.Err(), "time", e.Time, "eventID", e.ID)
		if err := am.bus.Alert.Send(a); err != nil {
			am.log.Warn("Cannot send alert", log.Error, err)
		}
	}
}

// checkEvent generates an alert if needed and returns nil if none is needed
func (am *AlertMgr) checkEvent(e *msg.SzenarioEvtMsg) *msg.AlertMsg {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.status.AddEvent(e)
	szStatusGroup := am.status.Get(e.Name)
	if szStatusGroup == nil {
		panic("szenario status cannot be nil")
	}
	lvl := szStatusGroup.Level()
	szState, found := am.basicStates[e.Name]
	if found {
		e.IncidentID = szState.GetIncidentID()
	}
	if lvl == status.OK {
		if found {
			diff := e.Time.Sub(szState.End)
			am.log.Info("Close incident", log.Szenario, e.Name, "start", szState.Start, "end", szState.End, "duration", diff)
			if diff > am.reopenTime {
				am.log.Info("Clear old incident", log.Szenario, e.Name, "state", szState)
				delete(am.basicStates, e.Name)
			}
			szState.End = e.Time
			i := msg.NewIncidentMsg(msg.CloseIncident, e)
			i.Start = szState.Start
			i.IntLevel = int(szStatusGroup.Level())
			i.ByteState = am.status.JSONBySzenario(e.Name)
			if err := core.Get().Bus().Incident.Send(i); err != nil {
				am.log.Warn("Cannot send incident", log.Error, err)
			}
		}
		return nil
	}
	if szState == nil {
		am.basicStates[e.Name] = newBasicState(am, e)
		i := msg.NewIncidentMsg(msg.OpenIncident, e)
		i.IntLevel = int(szStatusGroup.Level())
		i.ByteState = am.status.JSONBySzenario(e.Name)
		if err := core.Get().Bus().Incident.Send(i); err != nil {
			am.log.Warn("Cannot send incident", log.Error, err, log.Szenario, e.Name, "message", e.Err())
		}
		am.log.Debug("Not alerting: first alert", log.Szenario, e.Name, "message", e.Err())
		return nil
	}
	szState.End = time.Time{}
	if e.Err() != nil {
		i := msg.NewIncidentMsg(msg.UpdateIncident, e)
		i.Start = szState.Start
		i.IntLevel = int(lvl)
		i.ByteState = am.status.JSONBySzenario(e.Name)
		if err := core.Get().Bus().Incident.Send(i); err != nil {
			am.log.Warn("Cannot send incident", log.Error, err, log.Szenario, e.Name, "message", e.Err())
		}
	}
	if lvl < am.alertLevel {
		am.log.Info("NOT alerting level to low", log.Szenario, e.Name, "message", e.Err(), "min_level", lvl, "szenario_status", szStatusGroup)
		return nil
	}
	return szState.GetAlert(e, szStatusGroup)
}

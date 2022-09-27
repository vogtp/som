package alertmgr

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
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
	hcl            hcl.Logger
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
		hcl:            bus.GetLogger().Named("alertMgr"),
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
		am.hcl.Errorf("Cannot read state: %v", err)
	}
	am.Configure(options...)
	bus.Szenario.Handle(am.handle)
	am.hcl.Infof("AlertMgr started: alert level %q re-alerting every %v", am.alertLevel, am.alertIntervall)
}

// Configure the AlertMgr
func (am *AlertMgr) Configure(options ...Option) {
	for _, o := range options {
		o(am)
	}
}

func (am *AlertMgr) handle(e *msg.SzenarioEvtMsg) {
	defer func() {
		go am.save()
	}()
	am.hcl.Debugf("Got event from %v: %v", e.Name, e.Err())
	if a := am.checkEvent(e); a != nil {
		am.hcl.Warnf("Generating alert for %v: %v (%v, %v)", e.Name, e.Err(), e.Time, e.ID)
		am.bus.Alert.Send(a)
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
			am.hcl.Infof("Close incident: %s Start %v End %v Diff %v", e.Name, szState.Start, szState.End, diff)
			if diff > am.reopenTime {
				am.hcl.Infof("Clear old incident of %s: %+v", e.Name, szState)
				delete(am.basicStates, e.Name)
			}
			szState.End = e.Time
			i := msg.NewIncidentMsg(msg.CloseIncident, e)
			i.Start = szState.Start
			i.IntLevel = int(szStatusGroup.Level())
			i.ByteState = am.status.JSONBySzenario(e.Name)
			core.Get().Bus().Incident.Send(i)
		}
		return nil
	}
	if szState == nil {
		am.basicStates[e.Name] = newBasicState(am, e)
		i := msg.NewIncidentMsg(msg.OpenIncident, e)
		i.IntLevel = int(szStatusGroup.Level())
		i.ByteState = am.status.JSONBySzenario(e.Name)
		core.Get().Bus().Incident.Send(i)
		am.hcl.Debugf("Not alerting: first alert of %s: %v", e.Name, e.Err())
		return nil
	}
	szState.End = time.Time{}
	if e.Err() != nil {
		i := msg.NewIncidentMsg(msg.UpdateIncident, e)
		i.Start = szState.Start
		i.IntLevel = int(lvl)
		i.ByteState = am.status.JSONBySzenario(e.Name)
		core.Get().Bus().Incident.Send(i)
	}
	if lvl < am.alertLevel {
		am.hcl.Infof("Szenario %s alert level %v NOT alerting: %s", e.Name, lvl, szStatusGroup.StringInt(2))
		return nil
	}
	return szState.GetAlert(e, szStatusGroup)
}

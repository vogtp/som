package alertmgr

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

type basicState struct {
	am         *AlertMgr
	hcl        hcl.Logger
	Alerted    int
	err        error
	Start      time.Time
	End        time.Time
	Level      status.Level
	LastUpdate time.Time
	IncidentID string
}

func newBasicState(am *AlertMgr, e *msg.SzenarioEvtMsg) *basicState {
	s := &basicState{
		am:      am,
		hcl:     am.hcl.Named(e.Name),
		Alerted: 0,
		err:     e.Err(),
		Start:   e.Time,
	}
	e.IncidentID = s.GetIncidentID()
	return s
}

// GetIncidentID returns the incident ID or creates a new
func (s *basicState) GetIncidentID() string {
	if len(s.IncidentID) < 1 {
		s.IncidentID = uuid.NewString()
	}
	return s.IncidentID
}

// GetAlert returns the alert to be generated or nil if none is needed
// and updates bookkeeping
func (s *basicState) GetAlert(e *msg.SzenarioEvtMsg, statusGroup status.SzenarioGroup) *msg.AlertMsg {
	oldErr := s.err
	s.err = e.Err()
	if oldErr == nil || e.Err() == nil {
		// no alert the first time or alert was cleared
		s.hcl.Debugf("Not alerting not a constant error: %v was %v", e.Err(), oldErr)
		return nil
	}
	lvl := statusGroup.Level()
	if s.Alerted > 0 &&
		s.Level >= lvl &&
		s.LastUpdate.After(e.Time.Add(-1*s.am.alertIntervall)) {
		// alread alerted
		s.hcl.Debugf("Not alerting: %v -- last alert was %v (%s ago, not alerting more often than %v)", e.Err(), s.LastUpdate, e.Time.Sub(s.LastUpdate), s.am.alertIntervall)
		return nil
	}
	if s.Alerted > 1 {
		e.SetStatus("Realert", fmt.Sprintf("Alert is persistent %v alerts: re-alerting every %v", s.Alerted, s.am.alertIntervall))
	}
	if s.Alerted > 0 {
		e.Errors = append([]string{fmt.Sprintf("Initial error: %v (%s)", oldErr, s.LastUpdate.Format(cfg.TimeFormatString))}, e.Errors...)
	}
	if time.Since(s.LastUpdate) < viper.GetDuration(cfg.AlertDelay) {
		s.hcl.Infof("Not alerting %v alert is too young %v must be older than %v", e.Name, time.Since(e.Time), viper.GetDuration(cfg.AlertDelay))
		return nil
	}
	s.Level = lvl
	s.LastUpdate = e.Time
	s.Alerted++
	a := msg.NewAlert(e)
	a.SetStatus(KeyTopology, statusGroup.String())
	a.Level = lvl.String()
	s.hcl.Infof("Generating %s alert for %s", a.Level, a.Name)
	return a
}

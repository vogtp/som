package alertmgr

import (
	"fmt"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

type basicState struct {
	am         *AlertMgr
	log        *slog.Logger
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
		log:     am.log.With(log.Szenario, e.Name),
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
		s.log.Debug("Not alerting not a constant error", "message", e.Err(), "old_message", oldErr)
		return nil
	}
	lvl := statusGroup.Level()
	if s.Alerted > 0 &&
		s.Level >= lvl &&
		s.LastUpdate.After(e.Time.Add(-1*s.am.alertIntervall)) {
		// alread alerted
		s.log.Debug("Not alerting: not alerting too often", "message", e.Err(), "last", s.LastUpdate, "since", e.Time.Sub(s.LastUpdate), "min_intervall", s.am.alertIntervall)
		return nil
	}
	if s.Alerted > 1 {
		e.SetStatus("Realert", fmt.Sprintf("Alert is persistent %v alerts: re-alerting every %v", s.Alerted, s.am.alertIntervall))
	}
	if s.Alerted > 0 {
		e.Errors = append([]string{fmt.Sprintf("Initial error: %v (%s)", oldErr, s.LastUpdate.Format(cfg.TimeFormatString))}, e.Errors...)
	}
	if time.Since(s.LastUpdate) < viper.GetDuration(cfg.AlertDelay) {
		s.log.Info("Not alerting: alert is too young", "message", e.Name, "age", time.Since(e.Time), "min_age", viper.GetDuration(cfg.AlertDelay))
		return nil
	}
	s.Level = lvl
	s.LastUpdate = e.Time
	s.Alerted++
	a := msg.NewAlert(e)
	a.SetStatus(KeyTopology, statusGroup.String())
	a.Level = lvl.String()
	s.log.Info("Generating alert", "level", a.Level, log.Szenario, a.Name)
	return a
}

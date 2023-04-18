package webstatus

import (
	"context"
	"embed"
	"html/template"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"golang.org/x/exp/slog"
)

var (
	//go:embed templates static
	assetData embed.FS
	templates = template.Must(template.ParseFS(assetData, "templates/*.gohtml", "templates/common/*.gohtml"))
)

// WebStatus displays the current status on the web
type WebStatus struct {
	log      *slog.Logger
	data     *szenarioData
	dbAccess *db.Client
}

// New registers a WebStatus on the event bus
func New() *WebStatus {
	c := core.Get()
	s := &WebStatus{
		log: c.Log().With(log.Component, "webstatus"),
	}
	s.data = newSzenarioData(s.log)

	if err := s.data.load(); err != nil {
		s.log.Error("Cannot load config", log.Error, err)
	}
	c.Bus().Szenario.Handle(s.handleSzenarioEvt)
	c.Bus().Alert.Handle(s.handleAlert)
	c.Bus().Incident.Handle(s.handleIncident)
	s.routes()
	s.cleanup()
	return s
}

func (s *WebStatus) handleSzenarioEvt(e *msg.SzenarioEvtMsg) {
	s.log.Debug("Webstatus got event", log.Szenario, e.Name)
	s.data.mu.Lock()
	defer func() {
		s.data.mu.Unlock()
		go func() {
			if err := s.data.save(); err != nil {
				s.log.Error("Cannot save config", log.Error, err)
			}
		}()
	}()
	s.data.Status.AddEvent(e)
	s.data.Status.UpdatePrometheus()

	// calculate rolling average of availability
	for _, sz := range s.data.Status.Szenarios() {
		if sz.Key() != e.Name {
			continue
		}
		curAvail := sz.Availability()
		avail, found := s.data.Availabilites[e.Name]
		if !found {
			s.data.Availabilites[e.Name] = curAvail
			continue
		}
		s.data.Availabilites[e.Name] = (avail + curAvail) / 2
		s.log.Debug("Update availability", log.Szenario, e.Name, "old_availability", avail, "run_availability", curAvail, "new_availability", s.data.Availabilites[e.Name])
	}
}

func (s *WebStatus) handleAlert(a *msg.AlertMsg) {
	s.log.Debug("Webstatus got alert", log.Szenario, a.Name)
	if err := s.Ent().Alert.Save(context.Background(), a); err != nil {
		s.log.Error("Cannot save alert to DB", log.Error, err)
	}
}

func (s *WebStatus) handleIncident(i *msg.IncidentMsg) {
	s.log.Info("Webstatus got incident", "msg_type", i.Type.String(), log.Szenario, i.Name, "start", i.Start.Format(cfg.TimeFormatString), "end", i.End.Format(cfg.TimeFormatString))
	if err := s.Ent().Incident.Save(context.Background(), i); err != nil {
		s.log.Error("Cannot save incident to DB", log.Szenario, i.Name, "incident_id", i.ID, log.Error, err)
	}
}

// Ent returns the db.Access
func (s *WebStatus) Ent() *db.Client {
	if s.dbAccess == nil {
		entAccess, err := db.New()
		if err != nil {
			s.log.Error("Cannot connect to DB", log.Error, err)
		}
		s.dbAccess = entAccess
	}
	return s.dbAccess
}

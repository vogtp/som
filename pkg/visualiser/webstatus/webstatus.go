package webstatus

import (
	"context"
	"embed"
	"html/template"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

var (
	//go:embed templates static
	assetData embed.FS
	templates = template.Must(template.ParseFS(assetData, "templates/*.gohtml", "templates/common/*.gohtml"))
)

// WebStatus displays the current status on the web
type WebStatus struct {
	hcl      hcl.Logger
	data     *szenarioData
	dbAccess *db.Client
}

// New registers a WebStatus on the event bus
func New() *WebStatus {
	c := core.Get()
	s := &WebStatus{
		hcl: c.HCL().Named("webstatus"),
	}
	s.data = newSzenarioData(s.hcl)

	if err := s.data.load(); err != nil {
		s.hcl.Errorf("Cannot load config: %v", err)
	}
	c.Bus().Szenario.Handle(s.handleSzenarioEvt)
	c.Bus().Alert.Handle(s.handleAlert)
	c.Bus().Incident.Handle(s.handleIncident)
	s.routes()
	return s
}

func (s *WebStatus) getDataRoot() string {
	root := viper.GetString(cfg.DataDir)
	if len(root) > 0 && !strings.HasSuffix(root, "/") {
		root += "/"
	}
	return root
}

func (s *WebStatus) handleSzenarioEvt(e *msg.SzenarioEvtMsg) {
	s.hcl.Debugf("Webstatus got %s event", e.Name)
	s.data.mu.Lock()
	defer func() {
		s.data.mu.Unlock()
		go func() {
			if err := s.data.save(); err != nil {
				s.hcl.Errorf("Cannot save config: %v", err)
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
			avail = curAvail
			s.data.Availabilites[e.Name] = curAvail
			continue
		}
		s.data.Availabilites[e.Name] = (avail + curAvail) / 2
		s.hcl.Debugf("%s availability (%v + %v)/2 = %v", e.Name, avail, curAvail, s.data.Availabilites[e.Name])
	}
}

func (s *WebStatus) handleAlert(a *msg.AlertMsg) {
	s.hcl.Debugf("Webstatus got %s alert", a.Name)
	if err := s.Ent().Alert.Save(context.Background(), a); err != nil {
		s.hcl.Warnf("Cannot save alert to DB: %v", err)
	}
}

func (s *WebStatus) handleIncident(i *msg.IncidentMsg) {
	s.hcl.Infof("Webstatus got  %s %s (%s - %s) ", i.Type.String(), i.Name, i.Start.Format(cfg.TimeFormatString), i.End.Format(cfg.TimeFormatString))
	if err := s.Ent().Incident.Save(context.Background(), i); err != nil {
		s.hcl.Warnf("Cannot save incident to DB: %v", err)
	}
}

// Ent returns the db.Access
func (s *WebStatus) Ent() *db.Client {
	if s.dbAccess == nil {
		entAccess, err := db.New()
		if err != nil {
			s.hcl.Errorf("Cannot connect to DB: %v", err)
		}
		s.dbAccess = entAccess
	}

	return s.dbAccess
}

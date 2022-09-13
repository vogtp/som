package webstatus

import (
	"embed"
	"html/template"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database"
)

var (
	//go:embed templates static
	assetData embed.FS
	templates = template.Must(template.ParseFS(assetData, "templates/*.gohtml", "templates/common/*.gohtml"))
)

// WebStatus displays the current status on the web
type WebStatus struct {
	hcl              hcl.Logger
	data             *szenarioData
	alertFileTmpl    *template.Template
	alertPathTmpl    *template.Template
	incidentFileTmpl *template.Template
	incidentPathTmpl *template.Template
	muACache         sync.Mutex
	alertCache       map[string]string
	muICache         sync.Mutex
	incidentCache    map[string]string
	dbAccess         *database.Client
}

// New registers a WebStatus on the event bus
func New() *WebStatus {
	c := core.Get()
	s := &WebStatus{
		hcl:           c.HCL().Named("webstatus"),
		alertCache:    make(map[string]string),
		incidentCache: make(map[string]string),
	}
	s.initaliseAlertTemplates()
	s.initIncidentTemplates()
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
		avail, found := s.data.Availabilites[e.Name]
		if !found {
			s.data.Availabilites[e.Name] = sz.Availability()
			continue
		}
		s.data.Availabilites[e.Name] = (avail + sz.Availability()) / 2
		s.hcl.Debugf("%s availability (%v + %v)/2 = %v", e.Name, avail, sz.Availability(), s.data.Availabilites[e.Name])
	}

	//  do we need local timeseries?
	// s.handleTimeseries(e)
}

func (s *WebStatus) handleAlert(a *msg.AlertMsg) {
	s.hcl.Debugf("Webstatus got %s alert", a.Name)
	if err := s.saveAlert(a); err != nil {
		s.hcl.Warnf("Cannot save alert: %v", err)
	}
}

func (s *WebStatus) handleIncident(i *msg.IncidentMsg) {
	s.hcl.Infof("Webstatus got  %s %s (%s - %s) ", i.Type.String(), i.Name, i.Start.Format(cfg.TimeFormatString), i.End.Format(cfg.TimeFormatString))
	if err := s.saveIncident(i); err != nil {
		s.hcl.Warnf("Cannot save incident: %v", err)
	}
}

// Ent returns the db.Access
func (s *WebStatus) Ent() *database.Client {
	if s.dbAccess == nil {
		entAccess, err := database.New()
		if err != nil {
			s.hcl.Errorf("Cannot connect to DB: %v", err)
		}
		s.dbAccess = entAccess
	}

	return s.dbAccess
}

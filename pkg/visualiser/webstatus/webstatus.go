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
	"github.com/vogtp/som/pkg/visualiser/data"
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
}

// New registers a WebStatus on the event bus
func New() *WebStatus {
	c := core.Get()
	s := &WebStatus{
		hcl:           c.HCL().Named("webstatus"),
		alertCache:    make(map[string]string),
		incidentCache: make(map[string]string),
	}
	s.InitialertTemplates()
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

	//  do we need local timeseries?
	// s.handleTimeseries(e)
}

func (s *WebStatus) handleAlert(a *msg.AlertMsg) {
	s.hcl.Debugf("Webstatus got %s alert", a.Name)
	if err := s.saveAlert(a); err != nil {
		s.hcl.Warnf("Cannot save alert: %v", err)
	}
}

func (s *WebStatus) handleIncident(a *msg.IncidentMsg) {
	s.hcl.Debugf("Webstatus got %s incident %s", a.Name, a.Type.String())
	if err := s.saveIncident(a); err != nil {
		s.hcl.Warnf("Cannot save incident: %v", err)
	}
}

func (s *WebStatus) handleTimeseries(e *msg.SzenarioEvtMsg) {
	if e.Err() != nil {
		// do not record time in case of errors
		return
	}
	tot, ok := e.Counters["step.total"]
	if !ok {
		return
	}
	f, ok := tot.(float64)
	if !ok {
		return
	}
	ts, ok := s.data.Timeseries[e.Name]
	if !ok {
		ts = &data.Timeserie{}
		s.data.Timeseries[e.Name] = ts
	}
	ts.Add(e.Time, f)
}

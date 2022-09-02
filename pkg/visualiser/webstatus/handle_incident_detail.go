package webstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

const (
	// IncidentDetailPath is the path of the incitedent details
	IncidentDetailPath = "/incident/detail/"
)

type incidentData struct {
	db.IncidentModel
	Status   status.Status
	Errors   []db.ErrorModel
	Counters map[string]string
	Stati    map[string]string
	Files    []msg.FileMsgItem
	ErrStr   string
}

func (s *WebStatus) handleIncidentDetail(w http.ResponseWriter, r *http.Request) {
	id := ""
	idx := strings.Index(r.URL.Path, IncidentDetailPath)
	if idx < 1 {
		http.Error(w, "No incident ID given", http.StatusBadRequest)
		return
	}
	id = strings.ToLower(r.URL.Path[idx+len(IncidentDetailPath):])
	if strings.HasSuffix(id, "/") {
		id = id[:len(id)-1]
	}
	s.hcl.Debugf("incidents details %s requested", id)

	ctx := r.Context()
	a := s.DB()
	incidents, err := a.GetIncident(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aCnt := len(incidents)
	if aCnt < 1 {
		err = templates.ExecuteTemplate(w, "empty.gohtml", common("SOM No such Incident", r))
		if err != nil {
			s.hcl.Errorf("index Template error %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	var data = struct {
		*commonData
		PromURL    string
		Timeformat string
		FilesURL   string
		Name       string
		Start      time.Time
		End        time.Time
		Level      status.Level
		IncidentID string
		Incidents  []incidentData
		Alerts     []db.AlertModel
		AlertLink  string
	}{
		commonData: common("SOM Incident", r),
		IncidentID: id,
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Level:      status.Unknown,
		Timeformat: cfg.TimeFormatString,
		Incidents:  make([]incidentData, aCnt),
		Alerts:     make([]db.AlertModel, 0),
	}
	data.FilesURL = data.Baseurl + "/" + FilesPath
	data.AlertLink = data.Baseurl + "/" + AlertDetailPath
	s.hcl.Debugf("found %v incident records", aCnt)

	for i, f := range incidents {
		data.Name = f.Name
		data.Start = f.Start
		data.End = f.End

		if f.Level() > data.Level {
			data.Level = f.Level()
		}
		stat := status.New()
		err = json.Unmarshal(f.ByteState, stat)
		if err != nil {
			s.hcl.Warnf("Cannot unmarsh state of incident: %v", err)
		}

		id := incidentData{
			IncidentModel: f,
			Status:        prepaireStatus(stat),
			Files:         make([]msg.FileMsgItem, 0),
		}

		if alrts, err := a.GetAlertBy(ctx, "incident_id = ?", f.IncidentID); err == nil {
			data.Alerts = append(data.Alerts, alrts...)
		} else {
			s.hcl.Warnf("Loading alerts: %v", err)
		}

		id.ErrStr = id.Error
		if errs, err := a.GetErrors(ctx, f.ID); err == nil {
			id.Errors = errs
		} else {
			s.hcl.Warnf("Loading errors: %v", err)
		}
		if stati, err := a.GetStati(ctx, f.ID); err == nil {
			id.Stati = stati
		} else {
			s.hcl.Warnf("Loading stati: %v", err)
		}
		if ctrs, err := a.GetCounters(ctx, f.ID); err == nil {
			id.Counters = ctrs
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		if fils, err := a.GetFiles(ctx, f.ID); err == nil {
			id.Files = fils
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		data.Incidents[aCnt-i-1] = id
	}
	data.Title = fmt.Sprintf("SOM Incident: %s", data.Name)
	s.render(w, r, "incident_detail.gohtml", data)
}

package webstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

const (
	// IncidentDetailPath is the path of the incitedent details
	IncidentDetailPath = "/incident/detail/"
)

type incidentData struct {
	*ent.Incident
	Status   status.Status
	Errors   []*ent.Failure
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
	client := s.Ent()
	q := client.Incident.Query()
	if len(id) > 0 {
		uid, err := uuid.Parse(id)
		if err != nil {
			e := fmt.Errorf("cannot parse %s as uuid: %w", id, err)
			s.hcl.Error(e.Error())
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		q.Where(incident.IncidentID(uid))
	}
	incidents, err := q.All(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aCnt := len(incidents)
	if aCnt < 1 {
		err = templates.ExecuteTemplate(w, "empty.gohtml", common("SOM No such Incident", r))
		if err != nil {
			s.hcl.Errorf("incident details Template error %v", err)
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
		Alerts     []*ent.Alert
		AlertLink  string
	}{
		commonData: common("SOM Incident", r),
		IncidentID: id,
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Level:      status.Unknown,
		Timeformat: cfg.TimeFormatString,
		Incidents:  make([]incidentData, aCnt),
		Alerts:     make([]*ent.Alert, 0),
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
		err = json.Unmarshal(f.State, stat)
		if err != nil {
			s.hcl.Warnf("Cannot unmarsh state of incident: %v", err)
		}

		id := incidentData{
			Incident: f,
			Status:   prepaireStatus(stat),
			Stati:    make(map[string]string),
			Counters: make(map[string]string),
		}

		if alrts, err := client.Alert.Query().Where(alert.IncidentIDEQ(f.IncidentID)).All(ctx); err == nil {
			data.Alerts = append(data.Alerts, alrts...)
		} else {
			s.hcl.Warnf("Loading alerts: %v", err)
		}

		id.ErrStr = id.Error
		if errs, err := f.QueryFailures().All(ctx); err == nil {
			id.Errors = errs
		} else {
			s.hcl.Warnf("Loading errors: %v", err)
		}
		if stati, err := f.QueryStati().All(ctx); err == nil {
			for _, s := range stati {
				id.Stati[s.Name] = s.Value
			}
		} else {
			s.hcl.Warnf("Loading stati: %v", err)
		}
		if ctrs, err := f.QueryCounters().All(ctx); err == nil {
			for _, c := range ctrs {
				id.Counters[c.Name] = c.Value
			}
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		if fils, err := f.QueryFiles().All(ctx); err == nil {
			id.Files = make([]msg.FileMsgItem, len(fils))
			for i, f := range fils {
				id.Files[i] = f.MsgItem()
			}
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		data.Incidents[aCnt-i-1] = id
	}
	data.Title = fmt.Sprintf("SOM Incident: %s", data.Name)
	s.render(w, r, "incident_detail.gohtml", data)
}

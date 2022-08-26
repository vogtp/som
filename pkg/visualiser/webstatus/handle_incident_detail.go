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

	s.hcl.Infof("incidents details %s requested", id)
	//files, err := s.getIncidentDetailFiles(s.getIncidentRoot(), id)
	a := db.Access{}
	incidents, err := a.GetIncident(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aCnt := len(incidents)
	url := r.URL.String()
	if len(r.URL.RawQuery) > 0 {
		url += "&"
	} else {
		url += "?"
	}
	var data = struct {
		*commonData
		PromURL    string
		Timeformat string
		ThisURL    string
		Name       string
		Start      time.Time
		End        time.Time
		Level      status.Level
		IncidentID string
		Incidents  []incidentData
	}{
		commonData: common("SOM Incident", r),
		IncidentID: id,
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Level:      status.Unknown,
		Timeformat: cfg.TimeFormatString,
		ThisURL:    url,
		Incidents:  make([]incidentData, aCnt),
	}
	s.hcl.Infof("found %v incident recs", len(incidents))

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
			Counters:      make(map[string]string),
			Stati:         make(map[string]string),
			Files:         make([]msg.FileMsgItem, 0),
		}
		id.ErrStr = id.Error
		if errs, err := a.GetErrors(f.ID); err == nil {
			id.Errors = errs
		} else {
			s.hcl.Warnf("Loading errors: %v", err)
		}

		data.Incidents[aCnt-i-1] = id
	}
	data.Title = fmt.Sprintf("SOM Incident: %s", data.Name)

	for _, i := range data.Incidents {
		fmt.Printf("Err %v\n", i.Error)
	}

	// FIXME: migrate
	// r.ParseForm()
	// file := r.Form.Get("file")
	// if len(file) > 0 && len(data.Incidents) > 0 {
	// 	s.hcl.Infof("Serving file: %v", file)
	// 	parts := strings.Split(file, ".")
	// 	for _, f := range data.Incidents[0].Files {
	// 		if f.Name != parts[0] || f.Type.Ext != parts[1] {
	// 			continue
	// 		}
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Header().Add("Content-Type", f.Type.MimeType)
	// 		_, err := w.Write(f.Payload)
	// 		if err != nil {
	// 			s.hcl.Warnf("Cannot write file %s: %v", file, err)
	// 		}
	// 		return
	// 	}
	// 	return
	// }
	err = templates.ExecuteTemplate(w, "incident_detail.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

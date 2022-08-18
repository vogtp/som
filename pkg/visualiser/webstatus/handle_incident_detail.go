package webstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/status"
)

const (
	// IncidentDetailPath is the path of the incitedent details
	IncidentDetailPath = "/incident/detail/"
)

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
	files, err := s.getIncidentDetailFiles(s.getIncidentRoot(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aCnt := len(files)
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
		Incidents  []*incidentInfo
	}{
		commonData: common("SOM Incident", r),
		IncidentID: id,
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Level:      status.Unknown,
		Timeformat: cfg.TimeFormatString,
		ThisURL:    url,
		Incidents:  make([]*incidentInfo, aCnt),
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].IncidentInfo.Time.Before(files[j].IncidentInfo.Time)
	})
	for i, f := range files {
		data.Name = f.IncidentInfo.Name
		data.Start = f.IncidentInfo.Start
		data.End = f.IncidentInfo.End
		a, err := s.getIncident(f.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if a.Level > data.Level {
			data.Level = a.Level
		}
		stat := status.New()
		err = json.Unmarshal(a.ByteState, stat)
		if err != nil {
			s.hcl.Warnf("Cannot unmarsh state of incident: %v", err)
		}
		a.Status = prepaireStatus(stat)
		data.Incidents[aCnt-i-1] = a
	}
	data.Title = fmt.Sprintf("SOM Incident: %s", data.Name)

	r.ParseForm()
	file := r.Form.Get("file")
	if len(file) > 0 && len(data.Incidents) > 0 {
		s.hcl.Infof("Serving file: %v", file)
		parts := strings.Split(file, ".")
		for _, f := range data.Incidents[0].Files {
			if f.Name != parts[0] || f.Type.Ext != parts[1] {
				continue
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", f.Type.MimeType)
			_, err := w.Write(f.Payload)
			if err != nil {
				s.hcl.Warnf("Cannot write file %s: %v", file, err)
			}
			return
		}
		return
	}
	err = templates.ExecuteTemplate(w, "incident_detail.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

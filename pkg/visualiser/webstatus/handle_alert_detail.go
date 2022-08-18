package webstatus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
)

const (
	// AlertDetailPath is the path of the alert details
	AlertDetailPath = "/alert/detail/"
)

func (s *WebStatus) handleAlertDetail(w http.ResponseWriter, r *http.Request) {
	id := ""
	idx := strings.Index(r.URL.Path, AlertDetailPath)
	if idx < 1 {
		http.Error(w, "No alert ID given", http.StatusBadRequest)
		return
	}
	id = strings.ToLower(r.URL.Path[idx+len(AlertDetailPath):])
	if strings.HasSuffix(id, "/") {
		id = id[:len(id)-1]
	}

	s.hcl.Infof("alerts details %s requested", id)
	files, err := s.getAlertFiles(s.getAlertRoot(), id)
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
		PromURL            string
		Timeformat         string
		ThisURL            string
		IncidentDetailPath string
		Alerts             []*msg.AlertMsg
	}{
		commonData: common("SOM Alert Details", r),
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),

		Timeformat:         cfg.TimeFormatString,
		ThisURL:            url,
		IncidentDetailPath: IncidentDetailPath,
		Alerts:             make([]*msg.AlertMsg, aCnt),
	}
	for i, f := range files {
		a, err := s.getAlert(f.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.Alerts[aCnt-i-1] = a
	}
	r.ParseForm()
	file := r.Form.Get("file")
	if len(file) > 0 && len(data.Alerts) > 0 {
		s.hcl.Infof("Serving file: %v", file)
		parts := strings.Split(file, ".")
		for _, f := range data.Alerts[0].Files {
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
	err = templates.ExecuteTemplate(w, "alert_detail.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

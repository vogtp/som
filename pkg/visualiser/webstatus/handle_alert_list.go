package webstatus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

const (
	alertListPath = "/alert/list/"
)

func (s *WebStatus) handleAlertList(w http.ResponseWriter, r *http.Request) {
	sz := ""
	name := ""
	idx := strings.Index(r.URL.Path, alertListPath)
	if idx > 0 {
		name = r.URL.Path[idx+len(alertListPath):]
		sz = strings.ToLower(name)
		if strings.HasSuffix(sz, "/") {
			sz = sz[:len(sz)-1]
		}
	}
	if len(name) < 1 {
		name = "All Szenarios"
	}
	s.hcl.Infof("alerts for szenario %s requested", sz)
	files, err := s.getAlertFiles(s.getAlertRoot(), sz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data = struct {
		*commonData
		PromURL       string
		Timeformat    string
		AlertListPath string
		Alerts        []alertFile
		Szenarios     []string
		FilterName    string
	}{
		commonData:    common(fmt.Sprintf("SOM Alerts: %s", name), r),
		FilterName:    name,
		PromURL:       fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat:    cfg.TimeFormatString,
		AlertListPath: alertListPath,
		Alerts:        files,
	}
	for _, stat := range s.data.Status.Szenarios() {
		data.Szenarios = append(data.Szenarios, stat.Key())
	}
	err = templates.ExecuteTemplate(w, "alert_list.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

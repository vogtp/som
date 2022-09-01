package webstatus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

const (
	alertListPath = "/alert/list/"
)

type alertData struct {
	//Path       string
	//Name       string
	AlertInfo *db.AlertModel
	//Error      string
	DetailLink string
}

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
	ctx := r.Context()
	alerts, err := s.DB().GetAlertBySzenario(ctx, sz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	baseurl := core.Get().WebServer().BasePath()
	alertDatas := make([]alertData, len(alerts))
	for i, a := range alerts {
		alertDatas[i] = alertData{
			//Path:       a.ID.String(),
			//Name:       a.Name,
			AlertInfo: &a,
			//Error:      a.Error,
			DetailLink: fmt.Sprintf("%s/%s/%s/", baseurl, AlertDetailPath, a.ID),
		}
	}
	s.hcl.Infof("Loaded %v alerts", len(alerts))
	var data = struct {
		*commonData
		PromURL       string
		Timeformat    string
		AlertListPath string
		Alerts        []alertData
		Szenarios     []string
		FilterName    string
	}{
		commonData:    common(fmt.Sprintf("SOM Alerts: %s", name), r),
		FilterName:    name,
		PromURL:       fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat:    cfg.TimeFormatString,
		AlertListPath: alertListPath,
		Alerts:        alertDatas,
		Szenarios:     s.DB().AlertSzenarios(ctx),
	}
	err = templates.ExecuteTemplate(w, "alert_list.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

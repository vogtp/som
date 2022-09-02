package webstatus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

const (
	// AlertDetailPath is the path of the alert details
	AlertDetailPath = "/alert/detail/"
)

type alertDetailData struct {
	db.AlertModel
	Errors   []db.ErrorModel
	Counters map[string]string
	Stati    map[string]string
	Files    []msg.FileMsgItem
	ErrStr   string
}

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
	ctx := r.Context()
	alerts, err := s.DB().GetAlert(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aCnt := len(alerts)
	if aCnt < 1 {
		err = templates.ExecuteTemplate(w, "empty.gohtml", common("SOM No such Alert", r))
		if err != nil {
			s.hcl.Errorf("index Template error %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
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
		FilesURL           string
		IncidentDetailPath string
		Alerts             []alertDetailData
	}{
		commonData: common("SOM Alert Details", r),
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),

		Timeformat:         cfg.TimeFormatString,
		IncidentDetailPath: IncidentDetailPath,
		Alerts:             make([]alertDetailData, aCnt),
	}
	data.FilesURL = data.Baseurl + "/" + FilesPath

	for i, alert := range alerts {
		alertDetail := alertDetailData{
			AlertModel: alert,
		}
		alertDetail.ErrStr = alertDetail.Error
		if errs, err := s.DB().GetErrors(ctx, alert.ID); err == nil {
			alertDetail.Errors = errs
		} else {
			s.hcl.Warnf("Loading errors: %v", err)
		}
		if stati, err := s.DB().GetStati(ctx, alert.ID); err == nil {
			alertDetail.Stati = stati
		} else {
			s.hcl.Warnf("Loading stati: %v", err)
		}
		if ctrs, err := s.DB().GetCounters(ctx, alert.ID); err == nil {
			alertDetail.Counters = ctrs
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		if fils, err := s.DB().GetFiles(ctx, alert.ID); err == nil {
			alertDetail.Files = fils
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		data.Alerts[aCnt-i-1] = alertDetail
		s.hcl.Infof("Region: %v", alert.Region)
	}

	err = templates.ExecuteTemplate(w, "alert_detail.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

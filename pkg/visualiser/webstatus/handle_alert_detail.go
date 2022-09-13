package webstatus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent/alert"
)

const (
	// AlertDetailPath is the path of the alert details
	AlertDetailPath = "/alert/detail/"
)

type alertDetailData struct {
	*ent.Alert
	Errors   []*ent.Failure
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
	q := s.Ent().Alert.Query()
	if len(id) > 0 {
		uid, err := uuid.Parse(id)
		if err != nil {
			e := fmt.Errorf("cannot parse %s as uuid: %w", id, err)
			s.hcl.Error(e.Error())
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		q.Where(alert.UUID(uid))
	}
	alerts, err := q.All(ctx)
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
			Alert:    alert,
			Stati:    make(map[string]string),
			Counters: make(map[string]string),
		}
		alertDetail.ErrStr = alertDetail.Error
		if errs, err := alert.QueryFailures().All(ctx); err == nil {
			alertDetail.Errors = errs
		} else {
			s.hcl.Warnf("Loading errors: %v", err)
		}
		if stati, err := alert.QueryStati().All(ctx); err == nil {
			for _, s := range stati {
				alertDetail.Stati[s.Name] = s.Value
			}
		} else {
			s.hcl.Warnf("Loading stati: %v", err)
		}
		if ctrs, err := alert.QueryCounters().All(ctx); err == nil {
			for _, c := range ctrs {
				alertDetail.Counters[c.Name] = c.Value
			}
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		if fils, err := alert.QueryFiles().All(ctx); err == nil {
			alertDetail.Files = make([]msg.FileMsgItem, len(fils))
			for i, f := range fils {
				alertDetail.Files[i] = f.MsgItem()
			}
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		data.Alerts[aCnt-i-1] = alertDetail
		s.hcl.Infof("Region: %v", alert.Region)
	}
	s.render(w, r, "alert_detail.gohtml", data)
}

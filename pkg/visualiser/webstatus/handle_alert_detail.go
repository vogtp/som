package webstatus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/stater/alertmgr"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
)

const (
	// AlertDetailPath is the path of the alert details
	AlertDetailPath = "/alert/detail/"
)

type alertDetailData struct {
	*ent.Alert
	Errors   []*ent.Failure
	Counters map[string]float64
	Stati    map[string]string
	Steps    map[string]time.Duration
	Files    []msg.FileMsgItem
	ErrStr   string
}

func (s *WebStatus) handleAlertDetail(w http.ResponseWriter, r *http.Request) {
	id := ""
	idx := strings.Index(r.URL.Path, AlertDetailPath)
	if idx < 0 {
		s.Error(w, r, "No alert ID given", nil, http.StatusBadRequest)
		return
	}
	id = strings.ToLower(r.URL.Path[idx+len(AlertDetailPath):])
	id = strings.TrimSuffix(id, "/")

	s.log.Info("alerts details requested", "alert_id", id)
	ctx := r.Context()
	q := s.Ent().Alert.Query()
	if len(id) > 0 {
		uid, err := uuid.Parse(id)
		if err != nil {
			s.log.Error("cannot parse as uuid", "uuid", id, log.Error, err)
			s.Error(w, r, "Cannot parse UUID", err, http.StatusBadRequest)
			return
		}
		q.Where(alert.UUID(uid))
	}
	alerts, err := q.All(ctx)
	if err != nil {
		s.Error(w, r, "Database error alerts", err, http.StatusInternalServerError)
		return
	}

	aCnt := len(alerts)
	if aCnt < 1 {
		s.Error(w, r, "No such alert", err, http.StatusInternalServerError)
		return
	}
	// url := r.URL.String()
	// if len(r.URL.RawQuery) > 0 {
	// 	url += "&"
	// } else {
	// 	url += "?"
	// }
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
	data.FilesURL = data.Baseurl + FilesPath

	for i, alert := range alerts {
		if errors.Is(ctx.Err(), context.Canceled) {
			s.log.Info("Incident alert context canceld", log.Error, ctx.Err())
			return
		}
		alertDetail := alertDetailData{
			Alert:    alert,
			Stati:    make(map[string]string),
			Counters: make(map[string]float64),
			Steps:    make(map[string]time.Duration),
		}
		alertDetail.ErrStr = alertDetail.Error
		if errs, err := alert.QueryFailures().All(ctx); err == nil {
			alertDetail.Errors = errs
		} else {
			s.log.Warn("Failed loading errors", log.Error, err)
		}
		if stati, err := alert.QueryStati().All(ctx); err == nil {
			for _, s := range stati {
				if s.Name == alertmgr.KeyTopology {
					continue
				}
				alertDetail.Stati[s.Name] = s.Value
			}
		} else {
			s.log.Warn("Failed loading stati", log.Error, err)
		}
		const stepPrefix = "step."
		if ctrs, err := alert.QueryCounters().All(ctx); err == nil {
			for _, c := range ctrs {
				if strings.HasPrefix(c.Name, stepPrefix) {
					alertDetail.Steps[c.Name[len(stepPrefix):]] = time.Duration(c.Value * float64(time.Second)).Round(time.Millisecond)
					continue
				}
				alertDetail.Counters[c.Name] = c.Value
			}
		} else {
			s.log.Warn("Failed loading counters", log.Error, err)
		}
		if fils, err := alert.QueryFiles().All(ctx); err == nil {
			alertDetail.Files = make([]msg.FileMsgItem, len(fils))
			for i, f := range fils {
				alertDetail.Files[i] = f.MsgItem()
			}
		} else {
			s.log.Warn("Failed loading files", log.Error, err)
		}
		data.Alerts[aCnt-i-1] = alertDetail
		//s.log.Info("Region", "region", alert.Region)
	}
	s.render(w, r, "alert_detail.gohtml", data)
}

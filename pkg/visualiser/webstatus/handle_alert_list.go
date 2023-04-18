package webstatus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
)

const (
	alertListPath = "/alert/list/"
)

type alertData struct {
	AlertInfo  *ent.Alert
	DetailLink string
}

func (s *WebStatus) handleAlertList(w http.ResponseWriter, r *http.Request) {
	sz := ""
	name := ""
	idx := strings.Index(r.URL.Path, alertListPath)
	if idx > -1 {
		name = r.URL.Path[idx+len(alertListPath):]
		for strings.HasSuffix(name, "/") {
			name = name[:len(name)-1]
		}
		sz = strings.ToLower(name)
	}
	if len(name) < 1 {
		name = "All Szenarios"
	}
	s.log.Info("alerts  requested", "szenario", sz)
	common := common("SOM Alerts", r)
	ctx := r.Context()
	q := s.Ent().Alert.Query().Order(ent.Desc(alert.FieldTime))
	if len(sz) > 0 {
		s.log.Debug("Szenario filter", "szenario", sz)
		q.Where(alert.NameEqualFold(sz))
	}
	q.Where(
		alert.And(
			alert.TimeGTE(common.DatePicker.Start),
			alert.TimeLTE(common.DatePicker.End),
		),
	)
	alerts, err := q.All(ctx)
	if err != nil {
		s.Error(w, r, "No such alert list", err, http.StatusInternalServerError)
		return
	}
	baseurl := core.Get().WebServer().BasePath()
	alertDatas := make([]alertData, len(alerts))
	for i, a := range alerts {
		alertDatas[i] = alertData{
			AlertInfo:  a,
			DetailLink: fmt.Sprintf("%s%s/%s/", baseurl, AlertDetailPath, a.UUID.String()),
		}
	}
	s.log.Info("Loaded alerts", "size", len(alerts))

	szenarios, err := s.Ent().Alert.Szenarios(ctx)
	if err != nil {
		s.log.Warn("Cannot get list of szenarios", "error", err)
		if szenarios == nil {
			szenarios = make([]string, 0)
		}
	}

	common.Title = fmt.Sprintf("SOM Alerts: %s (%v)", name, len(alerts))
	var data = struct {
		*commonData
		PromURL       string
		Timeformat    string
		AlertListPath string
		Alerts        []alertData
		Szenarios     []string
		FilterName    string
	}{
		commonData:    common,
		FilterName:    name,
		PromURL:       fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat:    cfg.TimeFormatString,
		AlertListPath: alertListPath,
		Alerts:        alertDatas,
		Szenarios:     szenarios,
	}
	s.render(w, r, "alert_list.gohtml", data)
}

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
	incidentListPath = "/incident/list/"
)

func (s *WebStatus) handleIncidentList(w http.ResponseWriter, r *http.Request) {
	sz := ""
	name := ""
	idx := strings.Index(r.URL.Path, incidentListPath)
	if idx > 0 {
		name = r.URL.Path[idx+len(incidentListPath):]
		sz = strings.ToLower(name)
		if strings.HasSuffix(sz, "/") {
			sz = sz[:len(sz)-1]
		}
	}
	if len(name) < 1 {
		name = "All Szenarios"
	}
	s.hcl.Debugf("incidents for szenario %s requested", sz)

	ctx := r.Context()
	summary, err := s.DB().GetIncidentSummary(ctx, sz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	baseurl := core.Get().WebServer().BasePath()
	for i, s := range summary {
		summary[i].DetailLink = fmt.Sprintf("%s/%s/%s/", baseurl, IncidentDetailPath, s.IncidentID)
	}

	var data = struct {
		*commonData
		PromURL          string
		Timeformat       string
		IncidentListPath string
		Incidents        []db.IncidentSummary
		Szenarios        []string
	}{
		commonData:       common(fmt.Sprintf("SOM Incidents: %s", name), r),
		PromURL:          fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat:       cfg.TimeFormatString,
		IncidentListPath: incidentListPath,
		Incidents:        summary,
		Szenarios:        s.DB().IncidentSzenarios(ctx),
	}

	err = templates.ExecuteTemplate(w, "incident_list.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

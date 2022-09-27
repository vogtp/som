package webstatus

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

const (
	incidentListPath = "/incident/list/"
)

type incidentWeb struct {
	*db.IncidentSummary
	AlertCount int
}

func (s *WebStatus) handleIncidentList(w http.ResponseWriter, r *http.Request) {
	sz := ""
	name := ""
	idx := strings.Index(r.URL.Path, incidentListPath)
	if idx > 0 {
		name = r.URL.Path[idx+len(incidentListPath):]
		for strings.HasSuffix(name, "/") {
			name = name[:len(name)-1]
		}
		sz = strings.ToLower(name)
	}
	if len(name) < 1 {
		name = "All Szenarios"
	}
	s.hcl.Debugf("incidents for szenario %s requested", sz)
	common := common("SOM Incidents", r)
	ctx := r.Context()
	q := s.Ent().IncidentSummary.Query()
	if len(sz) > 0 {
		q.Where(incident.NameEqualFold(sz))
	}

	// q.Where(
	// 	incident.Not(
	// 		incident.Or(
	// 			incident.StartGT(common.End),
	// 			incident.And(
	// 				incident.EndNotNil(),
	// 				incident.EndLT(common.Start),
	// 			),
	// 		),
	// 	),
	// )

	summary, err := q.All(ctx)
	if err != nil {
		err = fmt.Errorf("Cannot load incidents from DB:\n %v", err)
		s.hcl.Errorf("Incident list: %v", err)
		s.Error(w, r, "Database error incident list", err, http.StatusInternalServerError)
		return
	}

	szenarios, err := s.Ent().Incident.Szenarios(ctx)
	if err != nil {
		s.hcl.Warnf("Cannot get list of szenarios: %v", err)
		if szenarios == nil {
			szenarios = make([]string, 0)
		}
	}

	common.Title = fmt.Sprintf("SOM Incidents: %s (%v)", name, len(summary))
	var data = struct {
		*commonData
		PromURL            string
		Timeformat         string
		IncidentListPath   string
		IncidentDetailPath string
		Incidents          []incidentWeb
		Szenarios          []string
		FilterName         string
	}{
		commonData:         common,
		FilterName:         name,
		PromURL:            fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat:         cfg.TimeFormatString,
		IncidentListPath:   incidentListPath,
		IncidentDetailPath: IncidentDetailPath,
		Szenarios:          szenarios,
	}
	filtered := make([]incidentWeb, 0)
	for _, sum := range summary {
		if sum.Start.Time().After(data.DatePicker.End) {
			continue
		}
		if !sum.End.IsZero() && sum.End.Time().Before(data.DatePicker.Start) {
			continue
		}
		sumWeb := incidentWeb{
			IncidentSummary: sum,
			AlertCount:      0,
		}
		if cnt, err := s.dbAccess.Alert.Query().Where(alert.IncidentIDEQ(sum.IncidentID)).Count(ctx); err == nil {
			sumWeb.AlertCount = cnt
		} else {
			s.hcl.Warnf("Cannot get alert count: %v", err)
		}
		filtered = append(filtered, sumWeb)
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].End.IsZero() && filtered[j].End.IsZero() {
			return filtered[i].Start.After(filtered[j].Start)
		}
		if filtered[i].End.IsZero() {
			return true
		}
		if filtered[j].End.IsZero() {
			return false
		}
		return filtered[i].Start.After(filtered[j].Start)
	})
	data.Incidents = filtered

	s.render(w, r, "incident_list.gohtml", data)
}

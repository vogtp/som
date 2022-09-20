package webstatus

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
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
		for strings.HasSuffix(name, "/") {
			name = name[:len(name)-1]
		}
		sz = strings.ToLower(name)
	}
	if len(name) < 1 {
		name = "All Szenarios"
	}
	s.hcl.Debugf("incidents for szenario %s requested", sz)

	ctx := r.Context()
	q := s.Ent().IncidentSummary.Query()
	// this works but the summaries are wrong
	//	q.Where(incident.And(incident.TimeGTE(start), incident.TimeLTE(end)))
	if len(sz) > 0 {
		q.Where(incident.NameEqualFold(sz))
	}

	summary, err := q.All(ctx)
	if err != nil {
		err = fmt.Errorf("Cannot load incidents from DB:\n %v", err)
		s.hcl.Errorf("Incident list: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	szenarios, err := s.Ent().Incident.Szenarios(ctx)
	if err != nil {
		s.hcl.Warnf("Cannot get list of szenarios: %v", err)
		if szenarios == nil {
			szenarios = make([]string, 0)
		}
	}

	var data = struct {
		*commonData
		PromURL            string
		Timeformat         string
		IncidentListPath   string
		IncidentDetailPath string
		Incidents          []*db.IncidentSummary
		Szenarios          []string
		FilterName         string
	}{
		commonData:         common(fmt.Sprintf("SOM Incidents: %s (%v)", name, len(summary)), r),
		FilterName:         name,
		PromURL:            fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat:         cfg.TimeFormatString,
		IncidentListPath:   incidentListPath,
		IncidentDetailPath: IncidentDetailPath,
		Szenarios:          szenarios,
	}
	filtered := make([]*db.IncidentSummary, 0)
	for _, s := range summary {
		if s.Start.Time().After(data.End) {
			continue
		}
		if !s.End.IsZero() && s.End.Time().Before(data.Start) {
			continue
		}

		filtered = append(filtered, s)
	}
	summary = filtered
	sort.Slice(summary, func(i, j int) bool {
		if summary[i].End.IsZero() && summary[j].End.IsZero() {
			return summary[i].Start.After(summary[j].Start)
		}
		if summary[i].End.IsZero() {
			return true
		}
		if summary[j].End.IsZero() {
			return false
		}
		return summary[i].Start.After(summary[j].Start)
	})
	data.Incidents = summary

	s.render(w, r, "incident_list.gohtml", data)
}

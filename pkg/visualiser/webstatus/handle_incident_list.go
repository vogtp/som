package webstatus

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
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
	s.hcl.Infof("incidents for szenario %s requested", sz)
	files, err := s.getIncidentFiles(s.getIncidentRoot(), sz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Slice(files, func(i, j int) bool {
		if (files[i].IncidentInfo.End.IsZero() || files[j].IncidentInfo.End.IsZero()) && !(files[i].IncidentInfo.End.IsZero() && files[j].IncidentInfo.End.IsZero()) {
			return files[i].IncidentInfo.End.IsZero() && !files[j].IncidentInfo.End.IsZero()
		}
		return files[i].IncidentInfo.Start.After(files[j].IncidentInfo.Start)
	})
	var data = struct {
		*commonData
		PromURL          string
		Timeformat       string
		IncidentListPath string
		Incidents        []incidentFile
		Szenarios        []string
	}{
		commonData:       common(fmt.Sprintf("SOM Incidents: %s", name), r),
		PromURL:          fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat:       cfg.TimeFormatString,
		IncidentListPath: incidentListPath,
		Incidents:        files,
	}

	for _, stat := range s.data.Status.Szenarios() {
		data.Szenarios = append(data.Szenarios, stat.Key())
	}

	err = templates.ExecuteTemplate(w, "incident_list.gohtml", data)
	if err != nil {
		s.hcl.Errorf("index Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

package webstatus

import (
	"fmt"
	"html/template"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/bridger"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

type indexValue struct {
	Name            string
	PromName        string
	Img             string
	Status          string
	AvailabilityCur string
	AvailabilityAvg string
	LastTime        template.HTML
	AvgTime         template.HTML
	LastUpdate      string
	IncidentCount   int
	IncidentList    string
}

func (s *WebStatus) handleIndex(w http.ResponseWriter, r *http.Request) {
	var data = struct {
		*commonData
		PromURL   string
		Szenarios []indexValue
	}{
		commonData: common("SOM Szenarios", r),
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
	}

	for _, stat := range s.data.Status.Szenarios() {
		szName := stat.Key()
		avail, found := s.data.Availabilites[szName]
		if !found {
			avail = stat.Availability()
		}

		s.hcl.Debugf("Displaying %s in index", szName)
		iv := indexValue{
			Name:            szName,
			PromName:        bridger.PrometheusName(szName),
			Status:          "\n" + stat.String(),
			AvailabilityCur: avail.String(),
			AvailabilityAvg: stat.Availability().String(),
			Img:             stat.Level().Img(),

			LastUpdate:   stat.LastUpdate().Format(cfg.TimeFormatString),
			IncidentList: fmt.Sprintf("%s%s/%s/", data.Baseurl, incidentListPath, szName),
		}

		iv.LastTime = formatStepTime(stat.LastTotal())
		times := stat.Totals()
		avg := 0.
		for _, t := range times {
			avg += t
		}
		avg /= float64(len(times))
		iv.AvgTime = formatStepTime(avg)
		var cnt = make([]struct {
			IncidentID uuid.UUID `json:"incident_id"`
		}, 0)
		if err := s.Ent().IncidentSummary.Query().Where(
			incident.NameEqualFold(szName),
			incident.And(
				incident.TimeGTE(data.DatePicker.Start),
				incident.TimeLTE(data.DatePicker.End),
			),
		).GroupBy(incident.FieldIncidentID).Scan(r.Context(), &cnt); err == nil {
			iv.IncidentCount = len(cnt)
		} else {
			iv.IncidentCount = -1
			s.hcl.Warnf("Cannot count incidents of %s: %v", szName, err)
		}
		data.Szenarios = append(data.Szenarios, iv)
	}
	s.render(w, r, "index.gohtml", data)
}

func formatStepTime(t float64) template.HTML {
	if math.IsNaN(t) {
		return ""
	}
	s := fmt.Sprintf("%.2fs", t)

	if t > 20 {
		s = fmt.Sprintf("<div style='color:red'>%s</div>", s)
	} else if t > 10 {
		s = fmt.Sprintf("<div style='color:orange'>%s</div>", s)
	}
	return template.HTML(s)

}

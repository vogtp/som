package webstatus

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/bridger"
	"github.com/vogtp/som/pkg/core/cfg"
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
		PromURL     string
		End         int
		Duration    int
		DurationStr string
		GraphStyle  string
		Szenarios   []indexValue
	}{
		commonData:  common("SOM Szenarios", r),
		PromURL:     fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Duration:    7 * 24 * 60 * 60,
		DurationStr: "7d",
		GraphStyle:  "avail",
	}
	s.hcl.Tracef("PromURL: %v", data.PromURL)

	utc := data.DatePicker.End.Hour() - data.DatePicker.End.UTC().Hour()
	data.End = int(data.DatePicker.End.Unix()) + utc*int(time.Hour.Seconds())
	data.Duration = int(data.DatePicker.End.Unix() - data.DatePicker.Start.Unix())

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
			IncidentList: fmt.Sprintf("%s/%s/%s/", data.Baseurl, incidentListPath, szName),
		}

		iv.LastTime = formatStepTime(stat.LastTotal())
		times := stat.Totals()
		avg := 0.
		for _, t := range times {
			avg += t
		}
		avg /= float64(len(times))
		iv.AvgTime = formatStepTime(avg)
		iv.IncidentCount = s.getIncidentCount(iv.Name)
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

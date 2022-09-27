package webstatus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

const (
	// IncidentDetailPath is the path of the incitedent details
	IncidentDetailPath = "/incident/detail/"
	pageSize           = 10
)

type incidentData struct {
	*ent.Incident
	Status   status.Status
	Errors   []*ent.Failure
	Counters map[string]string
	Stati    map[string]string
	Files    []msg.FileMsgItem
	ErrStr   string
}

type Page struct {
	ID    template.HTML
	State string
	URL   string
}

var start time.Time

func (s WebStatus) logTime(format string, v ...any) {
	s.hcl.Infof(format+" (%v)", append(v, time.Since(start))...)
}

func (s *WebStatus) handleIncidentDetail(w http.ResponseWriter, r *http.Request) {
	start = time.Now()
	id := ""
	idx := strings.Index(r.URL.Path, IncidentDetailPath)
	if idx < 1 {
		s.Error(w, r, "No incident ID given", nil, http.StatusBadRequest)
		return
	}
	id = strings.ToLower(r.URL.Path[idx+len(IncidentDetailPath):])
	if strings.HasSuffix(id, "/") {
		id = id[:len(id)-1]
	}
	s.hcl.Debugf("incidents details %s requested", id)

	ctx := r.Context()
	client := s.Ent()
	q := client.Incident.Query()
	var incidentID uuid.UUID
	if len(id) > 0 {
		var err error
		incidentID, err = uuid.Parse(id)
		if err != nil {
			e := fmt.Errorf("cannot parse %s as uuid: %w", id, err)
			s.hcl.Error(e.Error())
			s.Error(w, r, "Cannot parse UUID", e, http.StatusBadRequest)
			return
		}
		q.Where(incident.IncidentID(incidentID))
	}
	incidents, err := q.All(ctx)
	if err != nil {
		s.Error(w, r, "Database error incidents", err, http.StatusInternalServerError)
		return
	}
	incidentSummary, err := client.IncidentSummary.Query().Where(incident.IncidentIDEQ(incidentID)).First(ctx)
	if err != nil {
		s.Error(w, r, "Database error incident summaries", err, http.StatusInternalServerError)
		return
	}
	totalIncidents := len(incidents)
	aCnt := totalIncidents
	s.logTime("incident count: %v", totalIncidents)
	if aCnt < 1 {
		s.Error(w, r, "No such incident", err, http.StatusInternalServerError)
		return
	}
	var pages []Page
	if totalIncidents > pageSize {
		page := 1
		r.ParseForm()
		if str := r.Form.Get("page"); len(str) > 0 {
			if p, err := strconv.Atoi(str); err == nil {
				page = p
			} else {
				s.hcl.Warnf("Cannot parse page %q", p)
			}
		}
		offset := (page - 1) * pageSize
		s.logTime("Paging offset %v len %v total %v", offset, pageSize, totalIncidents)
		incidents, err = q.Offset(offset).Limit(pageSize).All(ctx)
		if err != nil {
			s.Error(w, r, "Database error incidents page", err, http.StatusInternalServerError)
			return
		}
		aCnt = len(incidents)
		pgCnt := int(math.Ceil(float64(totalIncidents / pageSize)))
		url := r.URL

		r.Form.Set("page", fmt.Sprintf("%d", page-1))
		r.URL.RawQuery = r.Form.Encode()
		p := Page{
			ID:  template.HTML("&laquo;"),
			URL: url.String(),
		}
		if page < 2 {
			p.State = "disabled"
		}
		pages = append(pages, p)
		for i := 1; i <= pgCnt; i++ {
			id := fmt.Sprintf("%v", i)
			r.Form.Set("page", id)
			r.URL.RawQuery = r.Form.Encode()
			p := Page{
				ID:  template.HTML(id),
				URL: url.String(),
			}
			if i == page {
				p.State = "active"
			}
			pages = append(pages, p)
		}
		if len(pages) > 18 {
			dots := Page{ID: "...", State: "disabled"}
			start := 9
			end := len(pages) - 9
			mid := []Page{dots}
			if !(page < start || page > end) {
				start -= 3
				end += 3
				mid = append(mid, pages[page-3:page+3]...)
				mid = append(mid, dots)
			}
			backP := pages[end:]
			pages = append(pages[:start], mid...)
			pages = append(pages, backP...)
		}
		r.Form.Set("page", fmt.Sprintf("%d", page+1))
		r.URL.RawQuery = r.Form.Encode()
		p = Page{
			ID:  template.HTML("&raquo;"),
			URL: url.String(),
		}
		if page >= pgCnt {
			p.State = "disabled"
		}
		pages = append(pages, p)
		r.Form.Del("page")
		r.URL.RawQuery = r.Form.Encode()
	}

	var data = struct {
		*commonData
		PromURL    string
		Timeformat string
		FilesURL   string
		Name       string
		Start      time.Time
		End        time.Time
		Level      status.Level
		IncidentID string
		Incidents  []incidentData
		Alerts     []*ent.Alert
		AlertLink  string
		Pages      []Page
	}{
		commonData: common("SOM Incident", r),
		IncidentID: id,
		Name:       incidentSummary.Name,
		Start:      incidentSummary.Start.Time(),
		End:        incidentSummary.End.Time(),
		Level:      incidentSummary.Level(),
		PromURL:    fmt.Sprintf("%v/%v", viper.GetString(cfg.PromURL), viper.GetString(cfg.PromBasePath)),
		Timeformat: cfg.TimeFormatString,
		Incidents:  make([]incidentData, aCnt),
		Alerts:     make([]*ent.Alert, 0),
		Pages:      pages,
	}
	data.TitleImage = fmt.Sprintf("%s/static/status/%s.png", data.Baseurl, data.Level.Img())
	data.FilesURL = data.Baseurl + "/" + FilesPath
	data.AlertLink = data.Baseurl + "/" + AlertDetailPath

	for i, f := range incidents {
		if errors.Is(ctx.Err(), context.Canceled) {
			s.hcl.Infof("Incident detail context canceld: %v", ctx.Err())
			return
		}

		if f.Level() > data.Level {
			data.Level = f.Level()
		}
		stat := status.New()
		err = json.Unmarshal(f.State, stat)
		if err != nil {
			s.hcl.Warnf("Cannot unmarsh state of incident: %v", err)
		}

		id := incidentData{
			Incident: f,
			Status:   prepaireStatus(stat),
			Stati:    make(map[string]string),
			Counters: make(map[string]string),
		}

		id.ErrStr = id.Error
		if errs, err := f.QueryFailures().All(ctx); err == nil {
			id.Errors = errs
		} else {
			s.hcl.Warnf("Loading errors: %v", err)
		}
		if stati, err := f.QueryStati().All(ctx); err == nil {
			for _, s := range stati {
				id.Stati[s.Name] = s.Value
			}
		} else {
			s.hcl.Warnf("Loading stati: %v", err)
		}
		if ctrs, err := f.QueryCounters().All(ctx); err == nil {
			for _, c := range ctrs {
				id.Counters[c.Name] = c.Value
			}
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		if fils, err := f.QueryFiles().Select(
			file.FieldUUID,
			file.FieldName,
			file.FieldType,
			file.FieldExt,
			file.FieldSize,
		).All(ctx); err == nil {
			id.Files = make([]msg.FileMsgItem, len(fils))
			for i, f := range fils {
				id.Files[i] = f.MsgItem()
			}
		} else {
			s.hcl.Warnf("Loading counters: %v", err)
		}
		data.Incidents[aCnt-i-1] = id
	}

	if alrts, err := client.Alert.Query().Where(alert.IncidentIDEQ(incidentID)).All(ctx); err == nil {
		data.Alerts = alrts
	} else {
		s.hcl.Warnf("Loading alerts: %v", err)
	}

	data.Title = fmt.Sprintf("SOM Incident: %s", data.Name)
	s.render(w, r, "incident_detail.gohtml", data)
}

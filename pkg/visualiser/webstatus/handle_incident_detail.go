package webstatus

import (
	"context"
	"encoding/json"
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
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/stater/alertmgr"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

const (
	// IncidentDetailPath is the path of the incitedent details
	IncidentDetailPath = "/incident/detail/"
)

type incidentData struct {
	*ent.Incident
	Status   status.Status
	Errors   []*ent.Failure
	Counters map[string]float64
	Stati    map[string]string
	Steps    map[string]time.Duration
	Files    []msg.FileMsgItem
	ErrStr   string
}

func (s *WebStatus) handleIncidentDetail(w http.ResponseWriter, r *http.Request) {
	id := ""
	idx := strings.Index(r.URL.Path, IncidentDetailPath)
	if idx < 0 {
		s.Error(w, r, "No incident ID given", nil, http.StatusBadRequest)
		return
	}
	id = strings.ToLower(r.URL.Path[idx+len(IncidentDetailPath):])
	id = strings.TrimSuffix(id, "/")
	s.log.Debug("incidents details requested", "id", id)

	ctx := r.Context()
	client := s.Ent()
	q := client.Incident.Query()
	var incidentID uuid.UUID
	if len(id) > 0 {
		var err error
		incidentID, err = uuid.Parse(id)
		if err != nil {
			e := fmt.Errorf("cannot parse %s as uuid: %w", id, err)
			s.log.Error("Cannot parse uuid", log.Error, e.Error(), "uuid", id)
			s.Error(w, r, "Cannot parse UUID", e, http.StatusBadRequest)
			return
		}
		q.Where(incident.IncidentID(incidentID))
	}

	incidentSummary, err := client.IncidentSummary.Query().Where(incident.IncidentIDEQ(incidentID)).First(ctx)
	if err != nil {
		s.Error(w, r, "Database error incident summaries", err, http.StatusInternalServerError)
		return
	}
	totalIncidents := incidentSummary.Total
	if totalIncidents < 1 {
		s.Error(w, r, "No such incident", err, http.StatusInternalServerError)
		return
	}

	pages, offset := s.getPages(r, totalIncidents)
	incidents, err := q.Offset(offset).Limit(pageSize).All(ctx)
	if err != nil {
		s.Error(w, r, "Database error incidents page", err, http.StatusInternalServerError)
		return
	}
	aCnt := len(incidents)

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
		Pages      []pageInfo
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
	data.FilesURL = data.Baseurl + FilesPath
	data.AlertLink = data.Baseurl + AlertDetailPath

	for i, f := range incidents {
		if errors.Is(ctx.Err(), context.Canceled) {
			s.log.Info("Incident detail context canceld", log.Error, ctx.Err())
			return
		}

		if f.Level() > data.Level {
			data.Level = f.Level()
		}
		stat := status.New()
		err = json.Unmarshal(f.State, stat)
		if err != nil {
			s.log.Warn("Cannot unmarsh state of incident", log.Error, err)
		}

		id := incidentData{
			Incident: f,
			Status:   prepaireStatus(stat),
			Stati:    make(map[string]string),
			Counters: make(map[string]float64),
			Steps:    make(map[string]time.Duration),
		}

		id.ErrStr = id.Error
		if errs, err := f.QueryFailures().All(ctx); err == nil {
			id.Errors = errs
		} else {
			s.log.Warn("Failed loading errors", log.Error, err)
		}
		if stati, err := f.QueryStati().All(ctx); err == nil {
			for _, s := range stati {
				if s.Name == alertmgr.KeyTopology {
					continue
				}
				id.Stati[s.Name] = s.Value
			}
		} else {
			s.log.Warn("Failed loading stati", log.Error, err)
		}
		const stepPrefix = "step."
		if ctrs, err := f.QueryCounters().All(ctx); err == nil {
			for _, c := range ctrs {
				if strings.HasPrefix(c.Name, stepPrefix) {
					id.Steps[c.Name[len(stepPrefix):]] = time.Duration(c.Value * float64(time.Second)).Round(time.Millisecond)
					continue
				}
				id.Counters[c.Name] = c.Value
			}
		} else {
			s.log.Warn("Failed loading counters", log.Error, err)
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
			s.log.Warn("Failed loading files", log.Error, err)
		}
		data.Incidents[aCnt-i-1] = id
	}

	if alrts, err := client.Alert.Query().Where(alert.IncidentIDEQ(incidentID)).All(ctx); err == nil {
		data.Alerts = alrts
	} else {
		s.log.Warn("Failed loading alerts: %v", log.Error, err)
	}

	data.Title = fmt.Sprintf("SOM Incident: %s", data.Name)
	s.render(w, r, "incident_detail.gohtml", data)
}

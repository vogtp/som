package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
)

// IncidentSummary is the summary of different incidents entries
type IncidentSummary struct {
	IncidentID uuid.UUID  `json:"incident_id"`
	Name       string     `json:"name"`
	Total      int        `json:"total"`
	IntLevel   int        `json:"level"`
	Start      MinMaxTime `json:"start"`
	End        MinMaxTime `json:"end"`
	LastUpdate MinMaxTime `json:"time"`
	Error      string     `json:"error"`
}

// IncidentSummaryQuery allows quering incidents
type IncidentSummaryQuery struct {
	client *Client
	*ent.IncidentQuery
}

// Query returns a list of incidents
func (isq *IncidentSummaryQuery) Query() *IncidentSummaryQuery {
	q := isq.client.Incident.Query()
	isq.IncidentQuery = q
	return isq
}

// Where wrap ent where
func (isq *IncidentSummaryQuery) Where(s ...predicate.Incident) *IncidentSummaryQuery {
	isq.IncidentQuery.Where(s...)
	return isq
}

// First wrap ent first
func (isq *IncidentSummaryQuery) First(ctx context.Context) (*IncidentSummary, error) {
	isq.IncidentQuery.Limit(1)
	all, err := isq.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, &ent.NotFoundError{}
	}
	return all[0], nil
}

// All returns all incudent summaries
func (isq *IncidentSummaryQuery) All(ctx context.Context) ([]*IncidentSummary, error) {
	var summary []*IncidentSummary
	err := isq.groupAndAggregate().Scan(ctx, &summary)
	// last event is the OK so remove it
	for i, s := range summary {
		summary[i].Total = s.Total - 1
	}
	return summary, err
}

func (isq *IncidentSummaryQuery) groupAndAggregate() *ent.IncidentGroupBy {
	return isq.IncidentQuery.Select(incident.FieldIncidentID, incident.FieldName).
		GroupBy(incident.FieldIncidentID, incident.FieldName).
		Aggregate(
			ent.As(ent.Count(), "Total"),
			ent.As(ent.Max(incident.FieldIntLevel), "Level"),
			ent.As(ent.Max(incident.FieldEnd), "End"),
			ent.As(ent.Max(incident.FieldTime), "Time"),
			ent.As(ent.Min(incident.FieldStart), "Start"),
			ent.As(ent.Max(incident.FieldError), "Error"),
		)
}

func (isq IncidentSummaryQuery) CloseIncident(ctx context.Context, is *IncidentSummary, username string, statusLevel status.Level, failure string) error {
	closer := isq.client.Incident.IncidentClient.Create()
	closer.SetUUID(uuid.New())
	closer.SetIncidentID(is.IncidentID)
	closer.SetName(is.Name)
	// fix stupid time formating
	now, _ := time.Parse(cfg.TimeFormatString, time.Now().Format(cfg.TimeFormatString))
	closer.SetTime(now)
	closer.SetStart(is.Start.Time())
	closer.SetEnd(now)
	closer.SetIntLevel(int(statusLevel))
	closer.SetError(failure)
	closer.SetUsername(username)
	closer.SetRegion("")
	closer.SetProbeHost("")
	closer.SetProbeOS("")
	closer.SetState([]byte(""))

	return closer.Exec(ctx)
}

// Level convinience method that calls status Level
func (is IncidentSummary) Level() status.Level {
	return status.Level(is.IntLevel)
}

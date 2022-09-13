package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent/incident"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

// IncidentSummary is the summary of different incidents entries
type IncidentSummary struct {
	IncidentID uuid.UUID     `json:"incident_id"`
	Name       string        `json:"name"`
	Total      int           `json:"total"`
	IntLevel   int           `json:"level"`
	Start      db.MinMaxTime `json:"start"`
	End        db.MinMaxTime `json:"end"`
	Error      string        `json:"error"`
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

// All returns all incudent summaries
func (isq *IncidentSummaryQuery) All(ctx context.Context) ([]*IncidentSummary, error) {
	g := isq.IncidentQuery.Select(incident.FieldIncidentID, incident.FieldName).
		GroupBy(incident.FieldIncidentID, incident.FieldName).
		Aggregate(
			ent.As(ent.Count(), "Total"),
			ent.As(ent.Max(incident.FieldIntLevel), "Level"),
			ent.As(ent.Max(incident.FieldEnd), "End"),
			ent.As(ent.Min(incident.FieldStart), "Start"),
			ent.As(ent.Max(incident.FieldError), "Error"),
		)
	var summary []*IncidentSummary
	err := g.Scan(ctx, &summary)
	return summary, err
}

// Level convinience method that calls status Level
func (is IncidentSummary) Level() status.Level {
	return status.Level(is.IntLevel)
}

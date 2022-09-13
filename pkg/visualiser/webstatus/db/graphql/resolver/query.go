package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
	graphql1 "github.com/vogtp/som/pkg/visualiser/webstatus/db/graphql"
)

// Incidents is the resolver for the Incidents field.
func (r *queryResolver) Incidents(ctx context.Context, szenario *string, timestamp *time.Time, incidentID *string, after *time.Time, before *time.Time) ([]*db.IncidentSummary, error) {
	q := r.client.IncidentSummary.Query()
	q.Order(ent.Desc(incident.FieldEnd))
	if incidentID != nil && len(*incidentID) > 0 {
		id, err := uuid.Parse(*incidentID)
		if err != nil {
			return nil, fmt.Errorf("%s is not a UUID: %w", *incidentID, err)
		}
		q.Where(incident.IncidentIDEQ(id))
	}
	if len(*szenario) > 0 {
		q.Where(incident.NameContains(*szenario))
	}
	if after != nil && !after.IsZero() {
		q.Where(incident.StartGTE(*after))
	}
	if before != nil && !before.IsZero() {
		q.Where(incident.And(incident.EndNEQ(time.Time{}), incident.EndLTE(*before)))
	}
	if timestamp != nil && !timestamp.IsZero() {
		q.Where(incident.And(incident.StartLTE(*timestamp), incident.EndGTE(*timestamp)))
	}
	return q.All(ctx)
}

// IncidentEntries is the resolver for the IncidentEntries field.
func (r *queryResolver) IncidentEntries(ctx context.Context, szenario *string, timestamp *time.Time, incidentID *string, after *time.Time, before *time.Time) ([]*ent.Incident, error) {
	q := r.client.Incident.Query().Order(ent.Desc(incident.FieldEnd))
	if incidentID != nil && len(*incidentID) > 0 {
		id, err := uuid.Parse(*incidentID)
		if err != nil {
			return nil, fmt.Errorf("%s is not a UUID: %w", *incidentID, err)
		}
		q.Where(incident.IncidentIDEQ(id))
	}
	if len(*szenario) > 0 {
		q.Where(incident.NameContains(*szenario))
	}
	if after != nil && !after.IsZero() {
		q.Where(incident.StartGTE(*after))
	}
	if before != nil && !before.IsZero() {
		q.Where(incident.And(incident.EndNEQ(time.Time{}), incident.EndLTE(*before)))
	}
	if timestamp != nil && !timestamp.IsZero() {
		q.Where(incident.And(incident.StartLTE(*timestamp), incident.EndGTE(*timestamp)))
	}
	return q.All(ctx)
}

// Alerts is the resolver for the Alerts field.
func (r *queryResolver) Alerts(ctx context.Context, szenario *string, after *time.Time, before *time.Time) ([]*ent.Alert, error) {
	q := r.client.Alert.Query().Order(ent.Desc(alert.FieldTime))
	if len(*szenario) > 0 {
		q.Where(alert.NameContains(*szenario))
	}
	if after != nil && !after.IsZero() {
		q.Where(alert.TimeGTE(*after))
	}
	if before != nil && !before.IsZero() {
		q.Where(alert.TimeLTE(*before))
	}
	return q.All(ctx)
}

// Query returns graphql1.QueryResolver implementation.
func (r *Resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
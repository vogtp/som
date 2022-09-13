package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/api/gqlgen"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

// Start is the resolver for the Start field.
func (r *incidentSummaryResolver) Start(ctx context.Context, obj *db.IncidentSummary) (*time.Time, error) {
	t := obj.Start.Time()
	return &t, nil
}

// End is the resolver for the End field.
func (r *incidentSummaryResolver) End(ctx context.Context, obj *db.IncidentSummary) (*time.Time, error) {
	t := obj.End.Time()
	return &t, nil
}

// IncidentEntries is the resolver for the IncidentEntries field.
func (r *incidentSummaryResolver) IncidentEntries(ctx context.Context, obj *db.IncidentSummary) ([]*ent.Incident, error) {
	return r.client.Incident.Query().Where(incident.IncidentID(obj.IncidentID)).All(ctx)
}

// Alerts is the resolver for the Alerts field.
func (r *incidentSummaryResolver) Alerts(ctx context.Context, obj *db.IncidentSummary, level *status.Level) ([]*ent.Alert, error) {
	q := r.client.Alert.Query().
		Order(ent.Desc(alert.FieldTime)).
		Where(alert.IncidentIDEQ(obj.IncidentID))
	if level != nil {
		lvl := *level
		if lvl > status.Unknown {
			q.Where(alert.IntLevelGTE(int(lvl)))
		}
	}
	return q.All(ctx)
}

// IncidentSummary returns gqlgen.IncidentSummaryResolver implementation.
func (r *Resolver) IncidentSummary() gqlgen.IncidentSummaryResolver {
	return &incidentSummaryResolver{r}
}

type incidentSummaryResolver struct{ *Resolver }

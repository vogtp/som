package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
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

// IncidentID is the resolver for the IncidentID field.
func (r *incidentSummaryResolver) IncidentID(ctx context.Context, obj *db.IncidentSummary) (string, error) {
	return obj.IncidentID.String(), nil
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

// IncidentSummary returns IncidentSummaryResolver implementation.
func (r *Resolver) IncidentSummary() IncidentSummaryResolver { return &incidentSummaryResolver{r} }

type incidentSummaryResolver struct{ *Resolver }

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
	graphql1 "github.com/vogtp/som/pkg/visualiser/webstatus/db/graphql"
)

// UUID is the resolver for the UUID field.
func (r *alertResolver) UUID(ctx context.Context, obj *ent.Alert) (string, error) {
	return obj.UUID.String(), nil
}

// IncidentID is the resolver for the IncidentID field.
func (r *alertResolver) IncidentID(ctx context.Context, obj *ent.Alert) (string, error) {
	return obj.IncidentID.String(), nil
}

// Incidents is the resolver for the Incidents field.
func (r *alertResolver) Incidents(ctx context.Context, obj *ent.Alert) ([]*db.IncidentSummary, error) {
	q := r.client.IncidentSummary.Query()
	q.Where(incident.IncidentIDEQ(obj.IncidentID))
	return q.All(ctx)
}

// IncidentEntries is the resolver for the IncidentEntries field.
func (r *alertResolver) IncidentEntries(ctx context.Context, obj *ent.Alert) ([]*ent.Incident, error) {
	return r.client.Incident.Query().Where(incident.IncidentID(obj.IncidentID)).All(ctx)
}

// Alert returns graphql1.AlertResolver implementation.
func (r *Resolver) Alert() graphql1.AlertResolver { return &alertResolver{r} }

type alertResolver struct{ *Resolver }

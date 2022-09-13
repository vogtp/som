package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
)

// State is the resolver for the State field.
func (r *incidentResolver) State(ctx context.Context, obj *ent.Incident) (string, error) {
	s := status.New()
	err := json.Unmarshal(obj.State, &s)
	if err != nil {
		return "", fmt.Errorf("Cannot unmarshal state: %w", err)
	}
	return s.String(), nil
}

// UUID is the resolver for the UUID field.
func (r *incidentResolver) UUID(ctx context.Context, obj *ent.Incident) (string, error) {
	return obj.UUID.String(), nil
}

// IncidentID is the resolver for the IncidentID field.
func (r *incidentResolver) IncidentID(ctx context.Context, obj *ent.Incident) (string, error) {
	return obj.IncidentID.String(), nil
}

// Alerts is the resolver for the Alerts field.
func (r *incidentResolver) Alerts(ctx context.Context, obj *ent.Incident) ([]*ent.Alert, error) {
	return r.client.Alert.Query().Where(alert.IncidentIDEQ(obj.IncidentID)).All(ctx)
}

// Incident returns IncidentResolver implementation.
func (r *Resolver) Incident() IncidentResolver { return &incidentResolver{r} }

type incidentResolver struct{ *Resolver }

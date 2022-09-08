package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/incident"
)

// Level is the resolver for the Level field.
func (r *incidentResolver) Level(ctx context.Context, obj *ent.Incident) (*string, error) {
	s := fmt.Sprintf("%v", obj.Level)
	return &s, nil
}

// State is the resolver for the State field.
func (r *incidentResolver) State(ctx context.Context, obj *ent.Incident) (*string, error) {
	panic(fmt.Errorf("not implemented: State - State"))
}

// UUID is the resolver for the UUID field.
func (r *incidentResolver) UUID(ctx context.Context, obj *ent.Incident) (*string, error) {
	s := (obj.UUID.String())
	return &s, nil
}

// IncidentID is the resolver for the IncidentID field.
func (r *incidentResolver) IncidentID(ctx context.Context, obj *ent.Incident) (*string, error) {
	s := obj.IncidentID.String()
	return &s, nil
}

// Incidents is the resolver for the Incidents field.
func (r *queryResolver) Incidents(ctx context.Context, szenario *string, timestamp *time.Time, start *time.Time, end *time.Time) ([]*ent.Incident, error) {
	q := r.client.Incident.Query().Order(ent.Desc(incident.FieldEnd))
	if len(*szenario) > 0 {
		q.Where(incident.NameContains(*szenario))
	}
	if start != nil && !start.IsZero() {
		hcl.Infof("Query start: %v", end)
		q.Where(incident.StartGTE(*start))
	}
	if end != nil && !end.IsZero() {
		hcl.Infof("Query end: %v", end)
		q.Where(incident.And(incident.EndNEQ(time.Time{}), incident.EndLTE(*end)))
	}
	if timestamp != nil && !timestamp.IsZero() {
		hcl.Infof("Query end: %v", end)
		q.Where(incident.And(incident.StartLTE(*timestamp), incident.EndGTE(*timestamp)))
	}
	return q.All(ctx)
}

// Incident returns IncidentResolver implementation.
func (r *Resolver) Incident() IncidentResolver { return &incidentResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type incidentResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/incident"
)

// Level is the resolver for the Level field.
func (r *alertResolver) Level(ctx context.Context, obj *ent.Alert) (*string, error) {
	s := fmt.Sprintf("%v", obj.Level)
	return &s, nil
}

// UUID is the resolver for the UUID field.
func (r *alertResolver) UUID(ctx context.Context, obj *ent.Alert) (*string, error) {
	s := obj.UUID.String()
	return &s, nil
}

// IncidentID is the resolver for the IncidentID field.
func (r *alertResolver) IncidentID(ctx context.Context, obj *ent.Alert) (*string, error) {
	s := obj.IncidentID.String()
	return &s, nil
}

// UUID is the resolver for the UUID field.
func (r *fileResolver) UUID(ctx context.Context, obj *ent.File) (*string, error) {
	s := obj.UUID.String()
	return &s, nil
}

// Payload is the resolver for the Payload field.
func (r *fileResolver) Payload(ctx context.Context, obj *ent.File) (*string, error) {
	s := string(obj.Payload)
	return &s, nil
}

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
func (r *queryResolver) Incidents(ctx context.Context, szenario *string, timestamp *time.Time, after *time.Time, before *time.Time) ([]*ent.Incident, error) {
	q := r.client.Incident.Query().Order(ent.Desc(incident.FieldEnd))
	if len(*szenario) > 0 {
		q.Where(incident.NameContains(*szenario))
	}
	if after != nil && !after.IsZero() {
		q.Where(incident.StartGTE(*after))
	}
	if before != nil && !before.IsZero() {
		hcl.Infof("Query end: %v", before)
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

// Alert returns AlertResolver implementation.
func (r *Resolver) Alert() AlertResolver { return &alertResolver{r} }

// File returns FileResolver implementation.
func (r *Resolver) File() FileResolver { return &fileResolver{r} }

// Incident returns IncidentResolver implementation.
func (r *Resolver) Incident() IncidentResolver { return &incidentResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type alertResolver struct{ *Resolver }
type fileResolver struct{ *Resolver }
type incidentResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

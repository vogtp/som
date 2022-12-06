package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/api/gqlgen"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

// Incident is the resolver for the Incident field.
func (r *subscriptionResolver) Incident(ctx context.Context, szenario *string) (<-chan *ent.Incident, error) {
	incChan := make(chan *ent.Incident)

	r.client.Incident.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate) {
				return next.Mutate(ctx, m)
			}
			val, errMut := next.Mutate(ctx, m)
			if szenario != nil && len(*szenario) > 0 {
				szVal, ok := m.Field(incident.FieldName)
				if !ok {
					return val, fmt.Errorf("Did not find szenario name")
				}
				szName, ok := szVal.(string)
				if !ok {
					return val, fmt.Errorf("Szenario name is not a string")
				}
				if szName != *szenario {
					return val, errMut
				}
			}
			uuidVal, ok := m.Field(incident.FieldUUID)
			if !ok {
				return val, errors.New("Mutation could not find uuid")
			}
			uuid, ok := uuidVal.(uuid.UUID)
			if !ok {
				return val, fmt.Errorf("cannot parse %v (%T) as uuid", uuidVal, uuidVal)
			}
			inc, err := r.client.Incident.Query().Where(incident.UUIDEQ(uuid)).First(ctx)
			if err != nil {
				return val, err
			}
			incChan <- inc
			return val, errMut
		})
	})

	return incChan, nil
}

// Alert is the resolver for the Alert field.
func (r *subscriptionResolver) Alert(ctx context.Context, szenario *string) (<-chan *ent.Alert, error) {
	alertChan := make(chan *ent.Alert)

	r.client.Alert.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !m.Op().Is(ent.OpCreate) {
				return next.Mutate(ctx, m)
			}
			val, errMut := next.Mutate(ctx, m)
			if szenario != nil && len(*szenario) > 0 {
				szVal, ok := m.Field(alert.FieldName)
				if !ok {
					return val, fmt.Errorf("Did not find szenario name")
				}
				szName, ok := szVal.(string)
				if !ok {
					return val, fmt.Errorf("Szenario name is not a string")
				}
				if szName != *szenario {
					return val, errMut
				}
			}
			uuidVal, ok := m.Field(alert.FieldUUID)
			if !ok {
				return val, errors.New("Mutation could not find uuid")
			}
			uuid, ok := uuidVal.(uuid.UUID)
			if !ok {
				return val, fmt.Errorf("cannot parse %v (%T) as uuid", uuidVal, uuidVal)
			}
			alert, err := r.client.Alert.Query().Where(alert.UUIDEQ(uuid)).First(ctx)
			if err != nil {
				return val, err
			}
			alertChan <- alert
			return val, errMut
		})
	})

	return alertChan, nil
}

// Subscription returns gqlgen.SubscriptionResolver implementation.
func (r *Resolver) Subscription() gqlgen.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }

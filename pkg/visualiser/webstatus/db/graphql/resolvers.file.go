package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
)

// UUID is the resolver for the UUID field.
func (r *fileResolver) UUID(ctx context.Context, obj *ent.File) (string, error) {
	return obj.UUID.String(), nil
}

// Payload is the resolver for the Payload field.
func (r *fileResolver) Payload(ctx context.Context, obj *ent.File) (string, error) {
	s := string(obj.Payload)
	// FIXME encode b64
	return s, nil
}

// File returns FileResolver implementation.
func (r *Resolver) File() FileResolver { return &fileResolver{r} }

type fileResolver struct{ *Resolver }

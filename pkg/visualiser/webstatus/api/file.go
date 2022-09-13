package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vogtp/som/pkg/visualiser/webstatus/api/gqlgen"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
)

// Payload is the resolver for the Payload field.
func (r *fileResolver) Payload(ctx context.Context, obj *ent.File) (string, error) {
	s := string(obj.Payload)
	// FIXME encode b64
	return s, nil
}

// File returns gqlgen.FileResolver implementation.
func (r *Resolver) File() gqlgen.FileResolver { return &fileResolver{r} }

type fileResolver struct{ *Resolver }

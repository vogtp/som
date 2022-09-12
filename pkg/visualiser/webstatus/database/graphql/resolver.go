package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver gives access to the DB
type Resolver struct {
	access *database.Access
	client *ent.Client
}

// NewSchema creates a graphql executable schema.
func NewSchema(a *database.Access) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{
			access: a,
			client: a.Client(),
		},
	})
}

package resolver

import (
	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver gives access to the DB
type Resolver struct {
	client *db.Client
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *db.Client) gqlgen.ExecutableSchema {
	return graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &Resolver{
			client: client,
		},
	})
}

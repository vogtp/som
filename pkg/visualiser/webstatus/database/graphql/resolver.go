package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver gives access to the DB
type Resolver struct {
	client *database.Client
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *database.Client) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{
			client: client,
		},
	})
}

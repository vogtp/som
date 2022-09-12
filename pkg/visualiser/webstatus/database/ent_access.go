package database

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent/migrate"
)

type Access struct {
	client *ent.Client
	hcl    hcl.Logger

	// IncidentSummary is the query for incident summaries
	IncidentSummary *IncidentSummaryQuery
}

// Client is the ent client
func (a Access) Client() *ent.Client {
	return a.client
}

// New creates an ent access
func New() (*Access, error) {
	client, err := ent.Open(dialect.SQLite, "file:data/somEnt.sqlite?&cache=shared&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}
	a := &Access{
		client: client,
		hcl:    core.Get().HCL().Named("ent"),
	}
	a.IncidentSummary = &IncidentSummaryQuery{a: a}

	return a, nil
}

// Close the client
func (a *Access) Close() {
	a.client.Close()
}

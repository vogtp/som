package database

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3" // needed to acces sqlite
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent/migrate"
)

// Client gives access to the DB and wraps ent
type Client struct {
	*ent.Client
	hcl hcl.Logger

	// IncidentSummary is the query for incident summaries
	IncidentSummary *IncidentSummaryQuery

	// Incident wraps and enhances the ent IncidentClient
	Incident *IncidentClient
	// Alert wraps and enhances the ent AlertClient
	Alert *AlertClient
}

// New creates an ent access
func New() (*Client, error) {
	entClient, err := ent.Open(dialect.SQLite, "file:data/somEnt.sqlite?&cache=shared&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := entClient.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}
	client := &Client{
		Client: entClient,
		hcl:    core.Get().HCL().Named("ent"),
	}
	client.IncidentSummary = &IncidentSummaryQuery{client: client}
	client.Incident = &IncidentClient{
		IncidentClient: entClient.Incident,
		client:         client,
	}
	client.Alert = &AlertClient{
		AlertClient: entClient.Alert,
		client:      client,
	}

	return client, nil
}

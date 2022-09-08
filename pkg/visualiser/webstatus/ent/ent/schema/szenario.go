package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

var (
	szFields = []ent.Field{
		field.UUID("UUID", uuid.UUID{}).Unique(),
		field.UUID("IncidentID", uuid.UUID{}),
		field.String("Name"),
		field.Time("Time"),
		field.String("Username"),
		field.String("Region"),
		field.String("ProbeOS"),
		field.String("ProbeHost"),
		field.String("Error").Optional(),
	}
	szEdges = []ent.Edge{
		edge.To("Counters", Counter.Type),
		edge.To("Stati", Status.Type),
		edge.To("Failures", Failure.Type),
	}
	szIndexes = []ent.Index{
		index.Fields("UUID").Unique().StorageKey("uuid"),
		index.Fields("IncidentID"),
		index.Fields("Name"),
	}
)

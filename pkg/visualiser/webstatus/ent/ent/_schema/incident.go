package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Incident holds the schema definition for the Incident entity.
type Incident struct {
	ent.Schema
}

// Fields of the Incident.
func (Incident) Fields() []ent.Field {
	return append([]ent.Field{
		field.Int("Level"),
		field.Time("Start"),
		field.Time("End"),
		field.Bytes("State"),
	}, szFields...)
}

// Edges of the Incident.
func (Incident) Edges() []ent.Edge {
	return szEdges
}

func (Incident) Indexes() []ent.Index {
	return append([]ent.Index{
		// non-unique index.
		index.Fields("Start"),
		index.Fields("End"),
		// unique index.
		// index.Fields("first_name", "last_name").
		//     Unique(),
	}, szIndexes...)
}

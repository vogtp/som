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

// Mixin of the Incident.
func (Incident) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SzenarioMixin{},
	}
}

// Fields of the Incident.
func (Incident) Fields() []ent.Field {
	return []ent.Field{
		field.Time("Start"),
		field.Time("End"),
		field.Bytes("State"),
	}
}

// Edges of the Incident.
func (Incident) Edges() []ent.Edge {
	return nil
}

// Indexes of the Incident.
func (Incident) Indexes() []ent.Index {
	return []ent.Index{
		// non-unique index.
		index.Fields("Start"),
		index.Fields("End"),
		// unique index.
		// index.Fields("first_name", "last_name").
		//     Unique(),
	}
}

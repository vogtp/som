package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Failure holds the schema definition for the Failure entity.
type Failure struct {
	ent.Schema
}

// Fields of the Err.
func (Failure) Fields() []ent.Field {
	return []ent.Field{
		field.String("Error"),
		field.Int("Idx"),
	}
}

// Edges of the Err.
func (Failure) Edges() []ent.Edge {
	return nil
}

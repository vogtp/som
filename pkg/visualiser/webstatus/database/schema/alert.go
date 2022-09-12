package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Alert holds the schema definition for the Alert entity.
type Alert struct {
	ent.Schema
}

// Fields of the Alert.
func (Alert) Fields() []ent.Field {
	return []ent.Field{
		field.Int("Level"),
	}
}

// Edges of the Alert.
func (Alert) Edges() []ent.Edge {
	return nil
}

// Mixin of the Alert.
func (Alert) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SzenarioMixin{},
	}
}
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// File holds the schema definition for the Counter entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("UUID", uuid.UUID{}),
		field.String("Name"),
		field.String("Type"),
		field.String("Ext"),
		field.Int("Size"),
		field.Bytes("payload"),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return nil
}

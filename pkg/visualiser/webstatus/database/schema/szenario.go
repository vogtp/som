package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// SzenarioMixin mixes common szenario fields in
type SzenarioMixin struct {
	mixin.Schema
}

// Fields of the SzenarioMixin.
func (SzenarioMixin) Fields() []ent.Field {
	return []ent.Field{
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
}

// Edges of the SzenarioMixin.
func (SzenarioMixin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Counters", Counter.Type).Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("Stati", Status.Type).Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("Failures", Failure.Type).Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("Files", File.Type).Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}

// Indexes of the SzenarioMixin.
func (SzenarioMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("UUID").Unique(),
		index.Fields("IncidentID"),
		index.Fields("Name"),
	}
}

// Code generated by ent, DO NOT EDIT.

package status

import (
	"entgo.io/ent/dialect/sql"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Status {
	return predicate.Status(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Status {
	return predicate.Status(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Status {
	return predicate.Status(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Status {
	return predicate.Status(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Status {
	return predicate.Status(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Status {
	return predicate.Status(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Status {
	return predicate.Status(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Status {
	return predicate.Status(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Status {
	return predicate.Status(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "Name" field. It's identical to NameEQ.
func Name(v string) predicate.Status {
	return predicate.Status(sql.FieldEQ(FieldName, v))
}

// Value applies equality check predicate on the "Value" field. It's identical to ValueEQ.
func Value(v string) predicate.Status {
	return predicate.Status(sql.FieldEQ(FieldValue, v))
}

// NameEQ applies the EQ predicate on the "Name" field.
func NameEQ(v string) predicate.Status {
	return predicate.Status(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "Name" field.
func NameNEQ(v string) predicate.Status {
	return predicate.Status(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "Name" field.
func NameIn(vs ...string) predicate.Status {
	return predicate.Status(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "Name" field.
func NameNotIn(vs ...string) predicate.Status {
	return predicate.Status(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "Name" field.
func NameGT(v string) predicate.Status {
	return predicate.Status(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "Name" field.
func NameGTE(v string) predicate.Status {
	return predicate.Status(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "Name" field.
func NameLT(v string) predicate.Status {
	return predicate.Status(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "Name" field.
func NameLTE(v string) predicate.Status {
	return predicate.Status(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "Name" field.
func NameContains(v string) predicate.Status {
	return predicate.Status(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "Name" field.
func NameHasPrefix(v string) predicate.Status {
	return predicate.Status(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "Name" field.
func NameHasSuffix(v string) predicate.Status {
	return predicate.Status(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "Name" field.
func NameEqualFold(v string) predicate.Status {
	return predicate.Status(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "Name" field.
func NameContainsFold(v string) predicate.Status {
	return predicate.Status(sql.FieldContainsFold(FieldName, v))
}

// ValueEQ applies the EQ predicate on the "Value" field.
func ValueEQ(v string) predicate.Status {
	return predicate.Status(sql.FieldEQ(FieldValue, v))
}

// ValueNEQ applies the NEQ predicate on the "Value" field.
func ValueNEQ(v string) predicate.Status {
	return predicate.Status(sql.FieldNEQ(FieldValue, v))
}

// ValueIn applies the In predicate on the "Value" field.
func ValueIn(vs ...string) predicate.Status {
	return predicate.Status(sql.FieldIn(FieldValue, vs...))
}

// ValueNotIn applies the NotIn predicate on the "Value" field.
func ValueNotIn(vs ...string) predicate.Status {
	return predicate.Status(sql.FieldNotIn(FieldValue, vs...))
}

// ValueGT applies the GT predicate on the "Value" field.
func ValueGT(v string) predicate.Status {
	return predicate.Status(sql.FieldGT(FieldValue, v))
}

// ValueGTE applies the GTE predicate on the "Value" field.
func ValueGTE(v string) predicate.Status {
	return predicate.Status(sql.FieldGTE(FieldValue, v))
}

// ValueLT applies the LT predicate on the "Value" field.
func ValueLT(v string) predicate.Status {
	return predicate.Status(sql.FieldLT(FieldValue, v))
}

// ValueLTE applies the LTE predicate on the "Value" field.
func ValueLTE(v string) predicate.Status {
	return predicate.Status(sql.FieldLTE(FieldValue, v))
}

// ValueContains applies the Contains predicate on the "Value" field.
func ValueContains(v string) predicate.Status {
	return predicate.Status(sql.FieldContains(FieldValue, v))
}

// ValueHasPrefix applies the HasPrefix predicate on the "Value" field.
func ValueHasPrefix(v string) predicate.Status {
	return predicate.Status(sql.FieldHasPrefix(FieldValue, v))
}

// ValueHasSuffix applies the HasSuffix predicate on the "Value" field.
func ValueHasSuffix(v string) predicate.Status {
	return predicate.Status(sql.FieldHasSuffix(FieldValue, v))
}

// ValueEqualFold applies the EqualFold predicate on the "Value" field.
func ValueEqualFold(v string) predicate.Status {
	return predicate.Status(sql.FieldEqualFold(FieldValue, v))
}

// ValueContainsFold applies the ContainsFold predicate on the "Value" field.
func ValueContainsFold(v string) predicate.Status {
	return predicate.Status(sql.FieldContainsFold(FieldValue, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Status) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Status) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Status) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		p(s.Not())
	})
}

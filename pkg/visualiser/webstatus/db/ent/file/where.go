// Code generated by ent, DO NOT EDIT.

package file

import (
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// UUID applies equality check predicate on the "UUID" field. It's identical to UUIDEQ.
func UUID(v uuid.UUID) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUUID), v))
	})
}

// Name applies equality check predicate on the "Name" field. It's identical to NameEQ.
func Name(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Type applies equality check predicate on the "Type" field. It's identical to TypeEQ.
func Type(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), v))
	})
}

// Ext applies equality check predicate on the "Ext" field. It's identical to ExtEQ.
func Ext(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExt), v))
	})
}

// Size applies equality check predicate on the "Size" field. It's identical to SizeEQ.
func Size(v int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSize), v))
	})
}

// Payload applies equality check predicate on the "payload" field. It's identical to PayloadEQ.
func Payload(v []byte) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPayload), v))
	})
}

// UUIDEQ applies the EQ predicate on the "UUID" field.
func UUIDEQ(v uuid.UUID) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUUID), v))
	})
}

// UUIDNEQ applies the NEQ predicate on the "UUID" field.
func UUIDNEQ(v uuid.UUID) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUUID), v))
	})
}

// UUIDIn applies the In predicate on the "UUID" field.
func UUIDIn(vs ...uuid.UUID) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldUUID), v...))
	})
}

// UUIDNotIn applies the NotIn predicate on the "UUID" field.
func UUIDNotIn(vs ...uuid.UUID) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldUUID), v...))
	})
}

// UUIDGT applies the GT predicate on the "UUID" field.
func UUIDGT(v uuid.UUID) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUUID), v))
	})
}

// UUIDGTE applies the GTE predicate on the "UUID" field.
func UUIDGTE(v uuid.UUID) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUUID), v))
	})
}

// UUIDLT applies the LT predicate on the "UUID" field.
func UUIDLT(v uuid.UUID) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUUID), v))
	})
}

// UUIDLTE applies the LTE predicate on the "UUID" field.
func UUIDLTE(v uuid.UUID) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUUID), v))
	})
}

// NameEQ applies the EQ predicate on the "Name" field.
func NameEQ(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "Name" field.
func NameNEQ(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "Name" field.
func NameIn(vs ...string) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "Name" field.
func NameNotIn(vs ...string) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "Name" field.
func NameGT(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "Name" field.
func NameGTE(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "Name" field.
func NameLT(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "Name" field.
func NameLTE(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "Name" field.
func NameContains(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "Name" field.
func NameHasPrefix(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "Name" field.
func NameHasSuffix(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "Name" field.
func NameEqualFold(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "Name" field.
func NameContainsFold(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// TypeEQ applies the EQ predicate on the "Type" field.
func TypeEQ(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), v))
	})
}

// TypeNEQ applies the NEQ predicate on the "Type" field.
func TypeNEQ(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldType), v))
	})
}

// TypeIn applies the In predicate on the "Type" field.
func TypeIn(vs ...string) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldType), v...))
	})
}

// TypeNotIn applies the NotIn predicate on the "Type" field.
func TypeNotIn(vs ...string) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldType), v...))
	})
}

// TypeGT applies the GT predicate on the "Type" field.
func TypeGT(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldType), v))
	})
}

// TypeGTE applies the GTE predicate on the "Type" field.
func TypeGTE(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldType), v))
	})
}

// TypeLT applies the LT predicate on the "Type" field.
func TypeLT(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldType), v))
	})
}

// TypeLTE applies the LTE predicate on the "Type" field.
func TypeLTE(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldType), v))
	})
}

// TypeContains applies the Contains predicate on the "Type" field.
func TypeContains(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldType), v))
	})
}

// TypeHasPrefix applies the HasPrefix predicate on the "Type" field.
func TypeHasPrefix(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldType), v))
	})
}

// TypeHasSuffix applies the HasSuffix predicate on the "Type" field.
func TypeHasSuffix(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldType), v))
	})
}

// TypeEqualFold applies the EqualFold predicate on the "Type" field.
func TypeEqualFold(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldType), v))
	})
}

// TypeContainsFold applies the ContainsFold predicate on the "Type" field.
func TypeContainsFold(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldType), v))
	})
}

// ExtEQ applies the EQ predicate on the "Ext" field.
func ExtEQ(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExt), v))
	})
}

// ExtNEQ applies the NEQ predicate on the "Ext" field.
func ExtNEQ(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldExt), v))
	})
}

// ExtIn applies the In predicate on the "Ext" field.
func ExtIn(vs ...string) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldExt), v...))
	})
}

// ExtNotIn applies the NotIn predicate on the "Ext" field.
func ExtNotIn(vs ...string) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldExt), v...))
	})
}

// ExtGT applies the GT predicate on the "Ext" field.
func ExtGT(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldExt), v))
	})
}

// ExtGTE applies the GTE predicate on the "Ext" field.
func ExtGTE(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldExt), v))
	})
}

// ExtLT applies the LT predicate on the "Ext" field.
func ExtLT(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldExt), v))
	})
}

// ExtLTE applies the LTE predicate on the "Ext" field.
func ExtLTE(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldExt), v))
	})
}

// ExtContains applies the Contains predicate on the "Ext" field.
func ExtContains(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldExt), v))
	})
}

// ExtHasPrefix applies the HasPrefix predicate on the "Ext" field.
func ExtHasPrefix(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldExt), v))
	})
}

// ExtHasSuffix applies the HasSuffix predicate on the "Ext" field.
func ExtHasSuffix(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldExt), v))
	})
}

// ExtEqualFold applies the EqualFold predicate on the "Ext" field.
func ExtEqualFold(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldExt), v))
	})
}

// ExtContainsFold applies the ContainsFold predicate on the "Ext" field.
func ExtContainsFold(v string) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldExt), v))
	})
}

// SizeEQ applies the EQ predicate on the "Size" field.
func SizeEQ(v int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSize), v))
	})
}

// SizeNEQ applies the NEQ predicate on the "Size" field.
func SizeNEQ(v int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSize), v))
	})
}

// SizeIn applies the In predicate on the "Size" field.
func SizeIn(vs ...int) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldSize), v...))
	})
}

// SizeNotIn applies the NotIn predicate on the "Size" field.
func SizeNotIn(vs ...int) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldSize), v...))
	})
}

// SizeGT applies the GT predicate on the "Size" field.
func SizeGT(v int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSize), v))
	})
}

// SizeGTE applies the GTE predicate on the "Size" field.
func SizeGTE(v int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSize), v))
	})
}

// SizeLT applies the LT predicate on the "Size" field.
func SizeLT(v int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSize), v))
	})
}

// SizeLTE applies the LTE predicate on the "Size" field.
func SizeLTE(v int) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSize), v))
	})
}

// PayloadEQ applies the EQ predicate on the "payload" field.
func PayloadEQ(v []byte) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPayload), v))
	})
}

// PayloadNEQ applies the NEQ predicate on the "payload" field.
func PayloadNEQ(v []byte) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPayload), v))
	})
}

// PayloadIn applies the In predicate on the "payload" field.
func PayloadIn(vs ...[]byte) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldPayload), v...))
	})
}

// PayloadNotIn applies the NotIn predicate on the "payload" field.
func PayloadNotIn(vs ...[]byte) predicate.File {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldPayload), v...))
	})
}

// PayloadGT applies the GT predicate on the "payload" field.
func PayloadGT(v []byte) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPayload), v))
	})
}

// PayloadGTE applies the GTE predicate on the "payload" field.
func PayloadGTE(v []byte) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPayload), v))
	})
}

// PayloadLT applies the LT predicate on the "payload" field.
func PayloadLT(v []byte) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPayload), v))
	})
}

// PayloadLTE applies the LTE predicate on the "payload" field.
func PayloadLTE(v []byte) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPayload), v))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.File) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.File) predicate.File {
	return predicate.File(func(s *sql.Selector) {
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
func Not(p predicate.File) predicate.File {
	return predicate.File(func(s *sql.Selector) {
		p(s.Not())
	})
}

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
)

// FileUpdate is the builder for updating File entities.
type FileUpdate struct {
	config
	hooks    []Hook
	mutation *FileMutation
}

// Where appends a list predicates to the FileUpdate builder.
func (fu *FileUpdate) Where(ps ...predicate.File) *FileUpdate {
	fu.mutation.Where(ps...)
	return fu
}

// SetUUID sets the "UUID" field.
func (fu *FileUpdate) SetUUID(u uuid.UUID) *FileUpdate {
	fu.mutation.SetUUID(u)
	return fu
}

// SetName sets the "Name" field.
func (fu *FileUpdate) SetName(s string) *FileUpdate {
	fu.mutation.SetName(s)
	return fu
}

// SetType sets the "Type" field.
func (fu *FileUpdate) SetType(s string) *FileUpdate {
	fu.mutation.SetType(s)
	return fu
}

// SetExt sets the "Ext" field.
func (fu *FileUpdate) SetExt(s string) *FileUpdate {
	fu.mutation.SetExt(s)
	return fu
}

// SetSize sets the "Size" field.
func (fu *FileUpdate) SetSize(i int) *FileUpdate {
	fu.mutation.ResetSize()
	fu.mutation.SetSize(i)
	return fu
}

// AddSize adds i to the "Size" field.
func (fu *FileUpdate) AddSize(i int) *FileUpdate {
	fu.mutation.AddSize(i)
	return fu
}

// SetPayload sets the "payload" field.
func (fu *FileUpdate) SetPayload(b []byte) *FileUpdate {
	fu.mutation.SetPayload(b)
	return fu
}

// Mutation returns the FileMutation object of the builder.
func (fu *FileUpdate) Mutation() *FileMutation {
	return fu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fu *FileUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, fu.sqlSave, fu.mutation, fu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fu *FileUpdate) SaveX(ctx context.Context) int {
	affected, err := fu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fu *FileUpdate) Exec(ctx context.Context) error {
	_, err := fu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fu *FileUpdate) ExecX(ctx context.Context) {
	if err := fu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fu *FileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(file.Table, file.Columns, sqlgraph.NewFieldSpec(file.FieldID, field.TypeInt))
	if ps := fu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fu.mutation.UUID(); ok {
		_spec.SetField(file.FieldUUID, field.TypeUUID, value)
	}
	if value, ok := fu.mutation.Name(); ok {
		_spec.SetField(file.FieldName, field.TypeString, value)
	}
	if value, ok := fu.mutation.GetType(); ok {
		_spec.SetField(file.FieldType, field.TypeString, value)
	}
	if value, ok := fu.mutation.Ext(); ok {
		_spec.SetField(file.FieldExt, field.TypeString, value)
	}
	if value, ok := fu.mutation.Size(); ok {
		_spec.SetField(file.FieldSize, field.TypeInt, value)
	}
	if value, ok := fu.mutation.AddedSize(); ok {
		_spec.AddField(file.FieldSize, field.TypeInt, value)
	}
	if value, ok := fu.mutation.Payload(); ok {
		_spec.SetField(file.FieldPayload, field.TypeBytes, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	fu.mutation.done = true
	return n, nil
}

// FileUpdateOne is the builder for updating a single File entity.
type FileUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FileMutation
}

// SetUUID sets the "UUID" field.
func (fuo *FileUpdateOne) SetUUID(u uuid.UUID) *FileUpdateOne {
	fuo.mutation.SetUUID(u)
	return fuo
}

// SetName sets the "Name" field.
func (fuo *FileUpdateOne) SetName(s string) *FileUpdateOne {
	fuo.mutation.SetName(s)
	return fuo
}

// SetType sets the "Type" field.
func (fuo *FileUpdateOne) SetType(s string) *FileUpdateOne {
	fuo.mutation.SetType(s)
	return fuo
}

// SetExt sets the "Ext" field.
func (fuo *FileUpdateOne) SetExt(s string) *FileUpdateOne {
	fuo.mutation.SetExt(s)
	return fuo
}

// SetSize sets the "Size" field.
func (fuo *FileUpdateOne) SetSize(i int) *FileUpdateOne {
	fuo.mutation.ResetSize()
	fuo.mutation.SetSize(i)
	return fuo
}

// AddSize adds i to the "Size" field.
func (fuo *FileUpdateOne) AddSize(i int) *FileUpdateOne {
	fuo.mutation.AddSize(i)
	return fuo
}

// SetPayload sets the "payload" field.
func (fuo *FileUpdateOne) SetPayload(b []byte) *FileUpdateOne {
	fuo.mutation.SetPayload(b)
	return fuo
}

// Mutation returns the FileMutation object of the builder.
func (fuo *FileUpdateOne) Mutation() *FileMutation {
	return fuo.mutation
}

// Where appends a list predicates to the FileUpdate builder.
func (fuo *FileUpdateOne) Where(ps ...predicate.File) *FileUpdateOne {
	fuo.mutation.Where(ps...)
	return fuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fuo *FileUpdateOne) Select(field string, fields ...string) *FileUpdateOne {
	fuo.fields = append([]string{field}, fields...)
	return fuo
}

// Save executes the query and returns the updated File entity.
func (fuo *FileUpdateOne) Save(ctx context.Context) (*File, error) {
	return withHooks(ctx, fuo.sqlSave, fuo.mutation, fuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fuo *FileUpdateOne) SaveX(ctx context.Context) *File {
	node, err := fuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fuo *FileUpdateOne) Exec(ctx context.Context) error {
	_, err := fuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuo *FileUpdateOne) ExecX(ctx context.Context) {
	if err := fuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fuo *FileUpdateOne) sqlSave(ctx context.Context) (_node *File, err error) {
	_spec := sqlgraph.NewUpdateSpec(file.Table, file.Columns, sqlgraph.NewFieldSpec(file.FieldID, field.TypeInt))
	id, ok := fuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "File.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, file.FieldID)
		for _, f := range fields {
			if !file.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != file.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fuo.mutation.UUID(); ok {
		_spec.SetField(file.FieldUUID, field.TypeUUID, value)
	}
	if value, ok := fuo.mutation.Name(); ok {
		_spec.SetField(file.FieldName, field.TypeString, value)
	}
	if value, ok := fuo.mutation.GetType(); ok {
		_spec.SetField(file.FieldType, field.TypeString, value)
	}
	if value, ok := fuo.mutation.Ext(); ok {
		_spec.SetField(file.FieldExt, field.TypeString, value)
	}
	if value, ok := fuo.mutation.Size(); ok {
		_spec.SetField(file.FieldSize, field.TypeInt, value)
	}
	if value, ok := fuo.mutation.AddedSize(); ok {
		_spec.AddField(file.FieldSize, field.TypeInt, value)
	}
	if value, ok := fuo.mutation.Payload(); ok {
		_spec.SetField(file.FieldPayload, field.TypeBytes, value)
	}
	_node = &File{config: fuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	fuo.mutation.done = true
	return _node, nil
}

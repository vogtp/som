// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
)

// FailureUpdate is the builder for updating Failure entities.
type FailureUpdate struct {
	config
	hooks    []Hook
	mutation *FailureMutation
}

// Where appends a list predicates to the FailureUpdate builder.
func (fu *FailureUpdate) Where(ps ...predicate.Failure) *FailureUpdate {
	fu.mutation.Where(ps...)
	return fu
}

// SetError sets the "Error" field.
func (fu *FailureUpdate) SetError(s string) *FailureUpdate {
	fu.mutation.SetError(s)
	return fu
}

// SetIdx sets the "Idx" field.
func (fu *FailureUpdate) SetIdx(i int) *FailureUpdate {
	fu.mutation.ResetIdx()
	fu.mutation.SetIdx(i)
	return fu
}

// AddIdx adds i to the "Idx" field.
func (fu *FailureUpdate) AddIdx(i int) *FailureUpdate {
	fu.mutation.AddIdx(i)
	return fu
}

// Mutation returns the FailureMutation object of the builder.
func (fu *FailureUpdate) Mutation() *FailureMutation {
	return fu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fu *FailureUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, FailureMutation](ctx, fu.sqlSave, fu.mutation, fu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fu *FailureUpdate) SaveX(ctx context.Context) int {
	affected, err := fu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fu *FailureUpdate) Exec(ctx context.Context) error {
	_, err := fu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fu *FailureUpdate) ExecX(ctx context.Context) {
	if err := fu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fu *FailureUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(failure.Table, failure.Columns, sqlgraph.NewFieldSpec(failure.FieldID, field.TypeInt))
	if ps := fu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fu.mutation.Error(); ok {
		_spec.SetField(failure.FieldError, field.TypeString, value)
	}
	if value, ok := fu.mutation.Idx(); ok {
		_spec.SetField(failure.FieldIdx, field.TypeInt, value)
	}
	if value, ok := fu.mutation.AddedIdx(); ok {
		_spec.AddField(failure.FieldIdx, field.TypeInt, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{failure.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	fu.mutation.done = true
	return n, nil
}

// FailureUpdateOne is the builder for updating a single Failure entity.
type FailureUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FailureMutation
}

// SetError sets the "Error" field.
func (fuo *FailureUpdateOne) SetError(s string) *FailureUpdateOne {
	fuo.mutation.SetError(s)
	return fuo
}

// SetIdx sets the "Idx" field.
func (fuo *FailureUpdateOne) SetIdx(i int) *FailureUpdateOne {
	fuo.mutation.ResetIdx()
	fuo.mutation.SetIdx(i)
	return fuo
}

// AddIdx adds i to the "Idx" field.
func (fuo *FailureUpdateOne) AddIdx(i int) *FailureUpdateOne {
	fuo.mutation.AddIdx(i)
	return fuo
}

// Mutation returns the FailureMutation object of the builder.
func (fuo *FailureUpdateOne) Mutation() *FailureMutation {
	return fuo.mutation
}

// Where appends a list predicates to the FailureUpdate builder.
func (fuo *FailureUpdateOne) Where(ps ...predicate.Failure) *FailureUpdateOne {
	fuo.mutation.Where(ps...)
	return fuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fuo *FailureUpdateOne) Select(field string, fields ...string) *FailureUpdateOne {
	fuo.fields = append([]string{field}, fields...)
	return fuo
}

// Save executes the query and returns the updated Failure entity.
func (fuo *FailureUpdateOne) Save(ctx context.Context) (*Failure, error) {
	return withHooks[*Failure, FailureMutation](ctx, fuo.sqlSave, fuo.mutation, fuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fuo *FailureUpdateOne) SaveX(ctx context.Context) *Failure {
	node, err := fuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fuo *FailureUpdateOne) Exec(ctx context.Context) error {
	_, err := fuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuo *FailureUpdateOne) ExecX(ctx context.Context) {
	if err := fuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fuo *FailureUpdateOne) sqlSave(ctx context.Context) (_node *Failure, err error) {
	_spec := sqlgraph.NewUpdateSpec(failure.Table, failure.Columns, sqlgraph.NewFieldSpec(failure.FieldID, field.TypeInt))
	id, ok := fuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Failure.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, failure.FieldID)
		for _, f := range fields {
			if !failure.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != failure.FieldID {
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
	if value, ok := fuo.mutation.Error(); ok {
		_spec.SetField(failure.FieldError, field.TypeString, value)
	}
	if value, ok := fuo.mutation.Idx(); ok {
		_spec.SetField(failure.FieldIdx, field.TypeInt, value)
	}
	if value, ok := fuo.mutation.AddedIdx(); ok {
		_spec.AddField(failure.FieldIdx, field.TypeInt, value)
	}
	_node = &Failure{config: fuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{failure.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	fuo.mutation.done = true
	return _node, nil
}

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent/failure"
)

// FailureCreate is the builder for creating a Failure entity.
type FailureCreate struct {
	config
	mutation *FailureMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetError sets the "Error" field.
func (fc *FailureCreate) SetError(s string) *FailureCreate {
	fc.mutation.SetError(s)
	return fc
}

// SetIdx sets the "Idx" field.
func (fc *FailureCreate) SetIdx(i int) *FailureCreate {
	fc.mutation.SetIdx(i)
	return fc
}

// Mutation returns the FailureMutation object of the builder.
func (fc *FailureCreate) Mutation() *FailureMutation {
	return fc.mutation
}

// Save creates the Failure in the database.
func (fc *FailureCreate) Save(ctx context.Context) (*Failure, error) {
	var (
		err  error
		node *Failure
	)
	if len(fc.hooks) == 0 {
		if err = fc.check(); err != nil {
			return nil, err
		}
		node, err = fc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FailureMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = fc.check(); err != nil {
				return nil, err
			}
			fc.mutation = mutation
			if node, err = fc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(fc.hooks) - 1; i >= 0; i-- {
			if fc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, fc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Failure)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from FailureMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FailureCreate) SaveX(ctx context.Context) *Failure {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fc *FailureCreate) Exec(ctx context.Context) error {
	_, err := fc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fc *FailureCreate) ExecX(ctx context.Context) {
	if err := fc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fc *FailureCreate) check() error {
	if _, ok := fc.mutation.Error(); !ok {
		return &ValidationError{Name: "Error", err: errors.New(`ent: missing required field "Failure.Error"`)}
	}
	if _, ok := fc.mutation.Idx(); !ok {
		return &ValidationError{Name: "Idx", err: errors.New(`ent: missing required field "Failure.Idx"`)}
	}
	return nil
}

func (fc *FailureCreate) sqlSave(ctx context.Context) (*Failure, error) {
	_node, _spec := fc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (fc *FailureCreate) createSpec() (*Failure, *sqlgraph.CreateSpec) {
	var (
		_node = &Failure{config: fc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: failure.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: failure.FieldID,
			},
		}
	)
	_spec.OnConflict = fc.conflict
	if value, ok := fc.mutation.Error(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: failure.FieldError,
		})
		_node.Error = value
	}
	if value, ok := fc.mutation.Idx(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: failure.FieldIdx,
		})
		_node.Idx = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Failure.Create().
//		SetError(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FailureUpsert) {
//			SetError(v+v).
//		}).
//		Exec(ctx)
func (fc *FailureCreate) OnConflict(opts ...sql.ConflictOption) *FailureUpsertOne {
	fc.conflict = opts
	return &FailureUpsertOne{
		create: fc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Failure.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fc *FailureCreate) OnConflictColumns(columns ...string) *FailureUpsertOne {
	fc.conflict = append(fc.conflict, sql.ConflictColumns(columns...))
	return &FailureUpsertOne{
		create: fc,
	}
}

type (
	// FailureUpsertOne is the builder for "upsert"-ing
	//  one Failure node.
	FailureUpsertOne struct {
		create *FailureCreate
	}

	// FailureUpsert is the "OnConflict" setter.
	FailureUpsert struct {
		*sql.UpdateSet
	}
)

// SetError sets the "Error" field.
func (u *FailureUpsert) SetError(v string) *FailureUpsert {
	u.Set(failure.FieldError, v)
	return u
}

// UpdateError sets the "Error" field to the value that was provided on create.
func (u *FailureUpsert) UpdateError() *FailureUpsert {
	u.SetExcluded(failure.FieldError)
	return u
}

// SetIdx sets the "Idx" field.
func (u *FailureUpsert) SetIdx(v int) *FailureUpsert {
	u.Set(failure.FieldIdx, v)
	return u
}

// UpdateIdx sets the "Idx" field to the value that was provided on create.
func (u *FailureUpsert) UpdateIdx() *FailureUpsert {
	u.SetExcluded(failure.FieldIdx)
	return u
}

// AddIdx adds v to the "Idx" field.
func (u *FailureUpsert) AddIdx(v int) *FailureUpsert {
	u.Add(failure.FieldIdx, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Failure.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *FailureUpsertOne) UpdateNewValues() *FailureUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Failure.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *FailureUpsertOne) Ignore() *FailureUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FailureUpsertOne) DoNothing() *FailureUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FailureCreate.OnConflict
// documentation for more info.
func (u *FailureUpsertOne) Update(set func(*FailureUpsert)) *FailureUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FailureUpsert{UpdateSet: update})
	}))
	return u
}

// SetError sets the "Error" field.
func (u *FailureUpsertOne) SetError(v string) *FailureUpsertOne {
	return u.Update(func(s *FailureUpsert) {
		s.SetError(v)
	})
}

// UpdateError sets the "Error" field to the value that was provided on create.
func (u *FailureUpsertOne) UpdateError() *FailureUpsertOne {
	return u.Update(func(s *FailureUpsert) {
		s.UpdateError()
	})
}

// SetIdx sets the "Idx" field.
func (u *FailureUpsertOne) SetIdx(v int) *FailureUpsertOne {
	return u.Update(func(s *FailureUpsert) {
		s.SetIdx(v)
	})
}

// AddIdx adds v to the "Idx" field.
func (u *FailureUpsertOne) AddIdx(v int) *FailureUpsertOne {
	return u.Update(func(s *FailureUpsert) {
		s.AddIdx(v)
	})
}

// UpdateIdx sets the "Idx" field to the value that was provided on create.
func (u *FailureUpsertOne) UpdateIdx() *FailureUpsertOne {
	return u.Update(func(s *FailureUpsert) {
		s.UpdateIdx()
	})
}

// Exec executes the query.
func (u *FailureUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for FailureCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FailureUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *FailureUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *FailureUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// FailureCreateBulk is the builder for creating many Failure entities in bulk.
type FailureCreateBulk struct {
	config
	builders []*FailureCreate
	conflict []sql.ConflictOption
}

// Save creates the Failure entities in the database.
func (fcb *FailureCreateBulk) Save(ctx context.Context) ([]*Failure, error) {
	specs := make([]*sqlgraph.CreateSpec, len(fcb.builders))
	nodes := make([]*Failure, len(fcb.builders))
	mutators := make([]Mutator, len(fcb.builders))
	for i := range fcb.builders {
		func(i int, root context.Context) {
			builder := fcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FailureMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, fcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = fcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, fcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fcb *FailureCreateBulk) SaveX(ctx context.Context) []*Failure {
	v, err := fcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fcb *FailureCreateBulk) Exec(ctx context.Context) error {
	_, err := fcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcb *FailureCreateBulk) ExecX(ctx context.Context) {
	if err := fcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Failure.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FailureUpsert) {
//			SetError(v+v).
//		}).
//		Exec(ctx)
func (fcb *FailureCreateBulk) OnConflict(opts ...sql.ConflictOption) *FailureUpsertBulk {
	fcb.conflict = opts
	return &FailureUpsertBulk{
		create: fcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Failure.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fcb *FailureCreateBulk) OnConflictColumns(columns ...string) *FailureUpsertBulk {
	fcb.conflict = append(fcb.conflict, sql.ConflictColumns(columns...))
	return &FailureUpsertBulk{
		create: fcb,
	}
}

// FailureUpsertBulk is the builder for "upsert"-ing
// a bulk of Failure nodes.
type FailureUpsertBulk struct {
	create *FailureCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Failure.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *FailureUpsertBulk) UpdateNewValues() *FailureUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Failure.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *FailureUpsertBulk) Ignore() *FailureUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FailureUpsertBulk) DoNothing() *FailureUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FailureCreateBulk.OnConflict
// documentation for more info.
func (u *FailureUpsertBulk) Update(set func(*FailureUpsert)) *FailureUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FailureUpsert{UpdateSet: update})
	}))
	return u
}

// SetError sets the "Error" field.
func (u *FailureUpsertBulk) SetError(v string) *FailureUpsertBulk {
	return u.Update(func(s *FailureUpsert) {
		s.SetError(v)
	})
}

// UpdateError sets the "Error" field to the value that was provided on create.
func (u *FailureUpsertBulk) UpdateError() *FailureUpsertBulk {
	return u.Update(func(s *FailureUpsert) {
		s.UpdateError()
	})
}

// SetIdx sets the "Idx" field.
func (u *FailureUpsertBulk) SetIdx(v int) *FailureUpsertBulk {
	return u.Update(func(s *FailureUpsert) {
		s.SetIdx(v)
	})
}

// AddIdx adds v to the "Idx" field.
func (u *FailureUpsertBulk) AddIdx(v int) *FailureUpsertBulk {
	return u.Update(func(s *FailureUpsert) {
		s.AddIdx(v)
	})
}

// UpdateIdx sets the "Idx" field to the value that was provided on create.
func (u *FailureUpsertBulk) UpdateIdx() *FailureUpsertBulk {
	return u.Update(func(s *FailureUpsert) {
		s.UpdateIdx()
	})
}

// Exec executes the query.
func (u *FailureUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the FailureCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for FailureCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FailureUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
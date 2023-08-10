// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/status"
)

// StatusCreate is the builder for creating a Status entity.
type StatusCreate struct {
	config
	mutation *StatusMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetName sets the "Name" field.
func (sc *StatusCreate) SetName(s string) *StatusCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetValue sets the "Value" field.
func (sc *StatusCreate) SetValue(s string) *StatusCreate {
	sc.mutation.SetValue(s)
	return sc
}

// Mutation returns the StatusMutation object of the builder.
func (sc *StatusCreate) Mutation() *StatusMutation {
	return sc.mutation
}

// Save creates the Status in the database.
func (sc *StatusCreate) Save(ctx context.Context) (*Status, error) {
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *StatusCreate) SaveX(ctx context.Context) *Status {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *StatusCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *StatusCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *StatusCreate) check() error {
	if _, ok := sc.mutation.Name(); !ok {
		return &ValidationError{Name: "Name", err: errors.New(`ent: missing required field "Status.Name"`)}
	}
	if _, ok := sc.mutation.Value(); !ok {
		return &ValidationError{Name: "Value", err: errors.New(`ent: missing required field "Status.Value"`)}
	}
	return nil
}

func (sc *StatusCreate) sqlSave(ctx context.Context) (*Status, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *StatusCreate) createSpec() (*Status, *sqlgraph.CreateSpec) {
	var (
		_node = &Status{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(status.Table, sqlgraph.NewFieldSpec(status.FieldID, field.TypeInt))
	)
	_spec.OnConflict = sc.conflict
	if value, ok := sc.mutation.Name(); ok {
		_spec.SetField(status.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sc.mutation.Value(); ok {
		_spec.SetField(status.FieldValue, field.TypeString, value)
		_node.Value = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Status.Create().
//		SetName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.StatusUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (sc *StatusCreate) OnConflict(opts ...sql.ConflictOption) *StatusUpsertOne {
	sc.conflict = opts
	return &StatusUpsertOne{
		create: sc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Status.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sc *StatusCreate) OnConflictColumns(columns ...string) *StatusUpsertOne {
	sc.conflict = append(sc.conflict, sql.ConflictColumns(columns...))
	return &StatusUpsertOne{
		create: sc,
	}
}

type (
	// StatusUpsertOne is the builder for "upsert"-ing
	//  one Status node.
	StatusUpsertOne struct {
		create *StatusCreate
	}

	// StatusUpsert is the "OnConflict" setter.
	StatusUpsert struct {
		*sql.UpdateSet
	}
)

// SetName sets the "Name" field.
func (u *StatusUpsert) SetName(v string) *StatusUpsert {
	u.Set(status.FieldName, v)
	return u
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *StatusUpsert) UpdateName() *StatusUpsert {
	u.SetExcluded(status.FieldName)
	return u
}

// SetValue sets the "Value" field.
func (u *StatusUpsert) SetValue(v string) *StatusUpsert {
	u.Set(status.FieldValue, v)
	return u
}

// UpdateValue sets the "Value" field to the value that was provided on create.
func (u *StatusUpsert) UpdateValue() *StatusUpsert {
	u.SetExcluded(status.FieldValue)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Status.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *StatusUpsertOne) UpdateNewValues() *StatusUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Status.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *StatusUpsertOne) Ignore() *StatusUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *StatusUpsertOne) DoNothing() *StatusUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the StatusCreate.OnConflict
// documentation for more info.
func (u *StatusUpsertOne) Update(set func(*StatusUpsert)) *StatusUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&StatusUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "Name" field.
func (u *StatusUpsertOne) SetName(v string) *StatusUpsertOne {
	return u.Update(func(s *StatusUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *StatusUpsertOne) UpdateName() *StatusUpsertOne {
	return u.Update(func(s *StatusUpsert) {
		s.UpdateName()
	})
}

// SetValue sets the "Value" field.
func (u *StatusUpsertOne) SetValue(v string) *StatusUpsertOne {
	return u.Update(func(s *StatusUpsert) {
		s.SetValue(v)
	})
}

// UpdateValue sets the "Value" field to the value that was provided on create.
func (u *StatusUpsertOne) UpdateValue() *StatusUpsertOne {
	return u.Update(func(s *StatusUpsert) {
		s.UpdateValue()
	})
}

// Exec executes the query.
func (u *StatusUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for StatusCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *StatusUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *StatusUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *StatusUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// StatusCreateBulk is the builder for creating many Status entities in bulk.
type StatusCreateBulk struct {
	config
	builders []*StatusCreate
	conflict []sql.ConflictOption
}

// Save creates the Status entities in the database.
func (scb *StatusCreateBulk) Save(ctx context.Context) ([]*Status, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Status, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*StatusMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = scb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *StatusCreateBulk) SaveX(ctx context.Context) []*Status {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *StatusCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *StatusCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Status.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.StatusUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (scb *StatusCreateBulk) OnConflict(opts ...sql.ConflictOption) *StatusUpsertBulk {
	scb.conflict = opts
	return &StatusUpsertBulk{
		create: scb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Status.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (scb *StatusCreateBulk) OnConflictColumns(columns ...string) *StatusUpsertBulk {
	scb.conflict = append(scb.conflict, sql.ConflictColumns(columns...))
	return &StatusUpsertBulk{
		create: scb,
	}
}

// StatusUpsertBulk is the builder for "upsert"-ing
// a bulk of Status nodes.
type StatusUpsertBulk struct {
	create *StatusCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Status.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *StatusUpsertBulk) UpdateNewValues() *StatusUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Status.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *StatusUpsertBulk) Ignore() *StatusUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *StatusUpsertBulk) DoNothing() *StatusUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the StatusCreateBulk.OnConflict
// documentation for more info.
func (u *StatusUpsertBulk) Update(set func(*StatusUpsert)) *StatusUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&StatusUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "Name" field.
func (u *StatusUpsertBulk) SetName(v string) *StatusUpsertBulk {
	return u.Update(func(s *StatusUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *StatusUpsertBulk) UpdateName() *StatusUpsertBulk {
	return u.Update(func(s *StatusUpsert) {
		s.UpdateName()
	})
}

// SetValue sets the "Value" field.
func (u *StatusUpsertBulk) SetValue(v string) *StatusUpsertBulk {
	return u.Update(func(s *StatusUpsert) {
		s.SetValue(v)
	})
}

// UpdateValue sets the "Value" field to the value that was provided on create.
func (u *StatusUpsertBulk) UpdateValue() *StatusUpsertBulk {
	return u.Update(func(s *StatusUpsert) {
		s.UpdateValue()
	})
}

// Exec executes the query.
func (u *StatusUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the StatusCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for StatusCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *StatusUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

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
)

// FileCreate is the builder for creating a File entity.
type FileCreate struct {
	config
	mutation *FileMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUUID sets the "UUID" field.
func (fc *FileCreate) SetUUID(u uuid.UUID) *FileCreate {
	fc.mutation.SetUUID(u)
	return fc
}

// SetName sets the "Name" field.
func (fc *FileCreate) SetName(s string) *FileCreate {
	fc.mutation.SetName(s)
	return fc
}

// SetType sets the "Type" field.
func (fc *FileCreate) SetType(s string) *FileCreate {
	fc.mutation.SetType(s)
	return fc
}

// SetExt sets the "Ext" field.
func (fc *FileCreate) SetExt(s string) *FileCreate {
	fc.mutation.SetExt(s)
	return fc
}

// SetSize sets the "Size" field.
func (fc *FileCreate) SetSize(i int) *FileCreate {
	fc.mutation.SetSize(i)
	return fc
}

// SetPayload sets the "payload" field.
func (fc *FileCreate) SetPayload(b []byte) *FileCreate {
	fc.mutation.SetPayload(b)
	return fc
}

// Mutation returns the FileMutation object of the builder.
func (fc *FileCreate) Mutation() *FileMutation {
	return fc.mutation
}

// Save creates the File in the database.
func (fc *FileCreate) Save(ctx context.Context) (*File, error) {
	var (
		err  error
		node *File
	)
	if len(fc.hooks) == 0 {
		if err = fc.check(); err != nil {
			return nil, err
		}
		node, err = fc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
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
		nv, ok := v.(*File)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from FileMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FileCreate) SaveX(ctx context.Context) *File {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fc *FileCreate) Exec(ctx context.Context) error {
	_, err := fc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fc *FileCreate) ExecX(ctx context.Context) {
	if err := fc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fc *FileCreate) check() error {
	if _, ok := fc.mutation.UUID(); !ok {
		return &ValidationError{Name: "UUID", err: errors.New(`ent: missing required field "File.UUID"`)}
	}
	if _, ok := fc.mutation.Name(); !ok {
		return &ValidationError{Name: "Name", err: errors.New(`ent: missing required field "File.Name"`)}
	}
	if _, ok := fc.mutation.GetType(); !ok {
		return &ValidationError{Name: "Type", err: errors.New(`ent: missing required field "File.Type"`)}
	}
	if _, ok := fc.mutation.Ext(); !ok {
		return &ValidationError{Name: "Ext", err: errors.New(`ent: missing required field "File.Ext"`)}
	}
	if _, ok := fc.mutation.Size(); !ok {
		return &ValidationError{Name: "Size", err: errors.New(`ent: missing required field "File.Size"`)}
	}
	if _, ok := fc.mutation.Payload(); !ok {
		return &ValidationError{Name: "payload", err: errors.New(`ent: missing required field "File.payload"`)}
	}
	return nil
}

func (fc *FileCreate) sqlSave(ctx context.Context) (*File, error) {
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

func (fc *FileCreate) createSpec() (*File, *sqlgraph.CreateSpec) {
	var (
		_node = &File{config: fc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: file.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		}
	)
	_spec.OnConflict = fc.conflict
	if value, ok := fc.mutation.UUID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: file.FieldUUID,
		})
		_node.UUID = value
	}
	if value, ok := fc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldName,
		})
		_node.Name = value
	}
	if value, ok := fc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
		_node.Type = value
	}
	if value, ok := fc.mutation.Ext(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldExt,
		})
		_node.Ext = value
	}
	if value, ok := fc.mutation.Size(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
		_node.Size = value
	}
	if value, ok := fc.mutation.Payload(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBytes,
			Value:  value,
			Column: file.FieldPayload,
		})
		_node.Payload = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.File.Create().
//		SetUUID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FileUpsert) {
//			SetUUID(v+v).
//		}).
//		Exec(ctx)
func (fc *FileCreate) OnConflict(opts ...sql.ConflictOption) *FileUpsertOne {
	fc.conflict = opts
	return &FileUpsertOne{
		create: fc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.File.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fc *FileCreate) OnConflictColumns(columns ...string) *FileUpsertOne {
	fc.conflict = append(fc.conflict, sql.ConflictColumns(columns...))
	return &FileUpsertOne{
		create: fc,
	}
}

type (
	// FileUpsertOne is the builder for "upsert"-ing
	//  one File node.
	FileUpsertOne struct {
		create *FileCreate
	}

	// FileUpsert is the "OnConflict" setter.
	FileUpsert struct {
		*sql.UpdateSet
	}
)

// SetUUID sets the "UUID" field.
func (u *FileUpsert) SetUUID(v uuid.UUID) *FileUpsert {
	u.Set(file.FieldUUID, v)
	return u
}

// UpdateUUID sets the "UUID" field to the value that was provided on create.
func (u *FileUpsert) UpdateUUID() *FileUpsert {
	u.SetExcluded(file.FieldUUID)
	return u
}

// SetName sets the "Name" field.
func (u *FileUpsert) SetName(v string) *FileUpsert {
	u.Set(file.FieldName, v)
	return u
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *FileUpsert) UpdateName() *FileUpsert {
	u.SetExcluded(file.FieldName)
	return u
}

// SetType sets the "Type" field.
func (u *FileUpsert) SetType(v string) *FileUpsert {
	u.Set(file.FieldType, v)
	return u
}

// UpdateType sets the "Type" field to the value that was provided on create.
func (u *FileUpsert) UpdateType() *FileUpsert {
	u.SetExcluded(file.FieldType)
	return u
}

// SetExt sets the "Ext" field.
func (u *FileUpsert) SetExt(v string) *FileUpsert {
	u.Set(file.FieldExt, v)
	return u
}

// UpdateExt sets the "Ext" field to the value that was provided on create.
func (u *FileUpsert) UpdateExt() *FileUpsert {
	u.SetExcluded(file.FieldExt)
	return u
}

// SetSize sets the "Size" field.
func (u *FileUpsert) SetSize(v int) *FileUpsert {
	u.Set(file.FieldSize, v)
	return u
}

// UpdateSize sets the "Size" field to the value that was provided on create.
func (u *FileUpsert) UpdateSize() *FileUpsert {
	u.SetExcluded(file.FieldSize)
	return u
}

// AddSize adds v to the "Size" field.
func (u *FileUpsert) AddSize(v int) *FileUpsert {
	u.Add(file.FieldSize, v)
	return u
}

// SetPayload sets the "payload" field.
func (u *FileUpsert) SetPayload(v []byte) *FileUpsert {
	u.Set(file.FieldPayload, v)
	return u
}

// UpdatePayload sets the "payload" field to the value that was provided on create.
func (u *FileUpsert) UpdatePayload() *FileUpsert {
	u.SetExcluded(file.FieldPayload)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.File.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *FileUpsertOne) UpdateNewValues() *FileUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.File.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *FileUpsertOne) Ignore() *FileUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FileUpsertOne) DoNothing() *FileUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FileCreate.OnConflict
// documentation for more info.
func (u *FileUpsertOne) Update(set func(*FileUpsert)) *FileUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FileUpsert{UpdateSet: update})
	}))
	return u
}

// SetUUID sets the "UUID" field.
func (u *FileUpsertOne) SetUUID(v uuid.UUID) *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.SetUUID(v)
	})
}

// UpdateUUID sets the "UUID" field to the value that was provided on create.
func (u *FileUpsertOne) UpdateUUID() *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.UpdateUUID()
	})
}

// SetName sets the "Name" field.
func (u *FileUpsertOne) SetName(v string) *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *FileUpsertOne) UpdateName() *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.UpdateName()
	})
}

// SetType sets the "Type" field.
func (u *FileUpsertOne) SetType(v string) *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "Type" field to the value that was provided on create.
func (u *FileUpsertOne) UpdateType() *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.UpdateType()
	})
}

// SetExt sets the "Ext" field.
func (u *FileUpsertOne) SetExt(v string) *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.SetExt(v)
	})
}

// UpdateExt sets the "Ext" field to the value that was provided on create.
func (u *FileUpsertOne) UpdateExt() *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.UpdateExt()
	})
}

// SetSize sets the "Size" field.
func (u *FileUpsertOne) SetSize(v int) *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.SetSize(v)
	})
}

// AddSize adds v to the "Size" field.
func (u *FileUpsertOne) AddSize(v int) *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.AddSize(v)
	})
}

// UpdateSize sets the "Size" field to the value that was provided on create.
func (u *FileUpsertOne) UpdateSize() *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.UpdateSize()
	})
}

// SetPayload sets the "payload" field.
func (u *FileUpsertOne) SetPayload(v []byte) *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.SetPayload(v)
	})
}

// UpdatePayload sets the "payload" field to the value that was provided on create.
func (u *FileUpsertOne) UpdatePayload() *FileUpsertOne {
	return u.Update(func(s *FileUpsert) {
		s.UpdatePayload()
	})
}

// Exec executes the query.
func (u *FileUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for FileCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FileUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *FileUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *FileUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// FileCreateBulk is the builder for creating many File entities in bulk.
type FileCreateBulk struct {
	config
	builders []*FileCreate
	conflict []sql.ConflictOption
}

// Save creates the File entities in the database.
func (fcb *FileCreateBulk) Save(ctx context.Context) ([]*File, error) {
	specs := make([]*sqlgraph.CreateSpec, len(fcb.builders))
	nodes := make([]*File, len(fcb.builders))
	mutators := make([]Mutator, len(fcb.builders))
	for i := range fcb.builders {
		func(i int, root context.Context) {
			builder := fcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FileMutation)
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
func (fcb *FileCreateBulk) SaveX(ctx context.Context) []*File {
	v, err := fcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fcb *FileCreateBulk) Exec(ctx context.Context) error {
	_, err := fcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcb *FileCreateBulk) ExecX(ctx context.Context) {
	if err := fcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.File.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FileUpsert) {
//			SetUUID(v+v).
//		}).
//		Exec(ctx)
func (fcb *FileCreateBulk) OnConflict(opts ...sql.ConflictOption) *FileUpsertBulk {
	fcb.conflict = opts
	return &FileUpsertBulk{
		create: fcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.File.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fcb *FileCreateBulk) OnConflictColumns(columns ...string) *FileUpsertBulk {
	fcb.conflict = append(fcb.conflict, sql.ConflictColumns(columns...))
	return &FileUpsertBulk{
		create: fcb,
	}
}

// FileUpsertBulk is the builder for "upsert"-ing
// a bulk of File nodes.
type FileUpsertBulk struct {
	create *FileCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.File.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *FileUpsertBulk) UpdateNewValues() *FileUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.File.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *FileUpsertBulk) Ignore() *FileUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FileUpsertBulk) DoNothing() *FileUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FileCreateBulk.OnConflict
// documentation for more info.
func (u *FileUpsertBulk) Update(set func(*FileUpsert)) *FileUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FileUpsert{UpdateSet: update})
	}))
	return u
}

// SetUUID sets the "UUID" field.
func (u *FileUpsertBulk) SetUUID(v uuid.UUID) *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.SetUUID(v)
	})
}

// UpdateUUID sets the "UUID" field to the value that was provided on create.
func (u *FileUpsertBulk) UpdateUUID() *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.UpdateUUID()
	})
}

// SetName sets the "Name" field.
func (u *FileUpsertBulk) SetName(v string) *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *FileUpsertBulk) UpdateName() *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.UpdateName()
	})
}

// SetType sets the "Type" field.
func (u *FileUpsertBulk) SetType(v string) *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "Type" field to the value that was provided on create.
func (u *FileUpsertBulk) UpdateType() *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.UpdateType()
	})
}

// SetExt sets the "Ext" field.
func (u *FileUpsertBulk) SetExt(v string) *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.SetExt(v)
	})
}

// UpdateExt sets the "Ext" field to the value that was provided on create.
func (u *FileUpsertBulk) UpdateExt() *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.UpdateExt()
	})
}

// SetSize sets the "Size" field.
func (u *FileUpsertBulk) SetSize(v int) *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.SetSize(v)
	})
}

// AddSize adds v to the "Size" field.
func (u *FileUpsertBulk) AddSize(v int) *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.AddSize(v)
	})
}

// UpdateSize sets the "Size" field to the value that was provided on create.
func (u *FileUpsertBulk) UpdateSize() *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.UpdateSize()
	})
}

// SetPayload sets the "payload" field.
func (u *FileUpsertBulk) SetPayload(v []byte) *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.SetPayload(v)
	})
}

// UpdatePayload sets the "payload" field to the value that was provided on create.
func (u *FileUpsertBulk) UpdatePayload() *FileUpsertBulk {
	return u.Update(func(s *FileUpsert) {
		s.UpdatePayload()
	})
}

// Exec executes the query.
func (u *FileUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the FileCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for FileCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FileUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
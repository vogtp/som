// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/counter"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/status"
)

// AlertCreate is the builder for creating a Alert entity.
type AlertCreate struct {
	config
	mutation *AlertMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUUID sets the "UUID" field.
func (ac *AlertCreate) SetUUID(u uuid.UUID) *AlertCreate {
	ac.mutation.SetUUID(u)
	return ac
}

// SetIncidentID sets the "IncidentID" field.
func (ac *AlertCreate) SetIncidentID(u uuid.UUID) *AlertCreate {
	ac.mutation.SetIncidentID(u)
	return ac
}

// SetName sets the "Name" field.
func (ac *AlertCreate) SetName(s string) *AlertCreate {
	ac.mutation.SetName(s)
	return ac
}

// SetTime sets the "Time" field.
func (ac *AlertCreate) SetTime(t time.Time) *AlertCreate {
	ac.mutation.SetTime(t)
	return ac
}

// SetIntLevel sets the "IntLevel" field.
func (ac *AlertCreate) SetIntLevel(i int) *AlertCreate {
	ac.mutation.SetIntLevel(i)
	return ac
}

// SetUsername sets the "Username" field.
func (ac *AlertCreate) SetUsername(s string) *AlertCreate {
	ac.mutation.SetUsername(s)
	return ac
}

// SetRegion sets the "Region" field.
func (ac *AlertCreate) SetRegion(s string) *AlertCreate {
	ac.mutation.SetRegion(s)
	return ac
}

// SetProbeOS sets the "ProbeOS" field.
func (ac *AlertCreate) SetProbeOS(s string) *AlertCreate {
	ac.mutation.SetProbeOS(s)
	return ac
}

// SetProbeHost sets the "ProbeHost" field.
func (ac *AlertCreate) SetProbeHost(s string) *AlertCreate {
	ac.mutation.SetProbeHost(s)
	return ac
}

// SetError sets the "Error" field.
func (ac *AlertCreate) SetError(s string) *AlertCreate {
	ac.mutation.SetError(s)
	return ac
}

// SetNillableError sets the "Error" field if the given value is not nil.
func (ac *AlertCreate) SetNillableError(s *string) *AlertCreate {
	if s != nil {
		ac.SetError(*s)
	}
	return ac
}

// AddCounterIDs adds the "Counters" edge to the Counter entity by IDs.
func (ac *AlertCreate) AddCounterIDs(ids ...int) *AlertCreate {
	ac.mutation.AddCounterIDs(ids...)
	return ac
}

// AddCounters adds the "Counters" edges to the Counter entity.
func (ac *AlertCreate) AddCounters(c ...*Counter) *AlertCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return ac.AddCounterIDs(ids...)
}

// AddStatiIDs adds the "Stati" edge to the Status entity by IDs.
func (ac *AlertCreate) AddStatiIDs(ids ...int) *AlertCreate {
	ac.mutation.AddStatiIDs(ids...)
	return ac
}

// AddStati adds the "Stati" edges to the Status entity.
func (ac *AlertCreate) AddStati(s ...*Status) *AlertCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ac.AddStatiIDs(ids...)
}

// AddFailureIDs adds the "Failures" edge to the Failure entity by IDs.
func (ac *AlertCreate) AddFailureIDs(ids ...int) *AlertCreate {
	ac.mutation.AddFailureIDs(ids...)
	return ac
}

// AddFailures adds the "Failures" edges to the Failure entity.
func (ac *AlertCreate) AddFailures(f ...*Failure) *AlertCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ac.AddFailureIDs(ids...)
}

// AddFileIDs adds the "Files" edge to the File entity by IDs.
func (ac *AlertCreate) AddFileIDs(ids ...int) *AlertCreate {
	ac.mutation.AddFileIDs(ids...)
	return ac
}

// AddFiles adds the "Files" edges to the File entity.
func (ac *AlertCreate) AddFiles(f ...*File) *AlertCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ac.AddFileIDs(ids...)
}

// Mutation returns the AlertMutation object of the builder.
func (ac *AlertCreate) Mutation() *AlertMutation {
	return ac.mutation
}

// Save creates the Alert in the database.
func (ac *AlertCreate) Save(ctx context.Context) (*Alert, error) {
	var (
		err  error
		node *Alert
	)
	if len(ac.hooks) == 0 {
		if err = ac.check(); err != nil {
			return nil, err
		}
		node, err = ac.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AlertMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ac.check(); err != nil {
				return nil, err
			}
			ac.mutation = mutation
			if node, err = ac.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ac.hooks) - 1; i >= 0; i-- {
			if ac.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ac.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, ac.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Alert)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from AlertMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ac *AlertCreate) SaveX(ctx context.Context) *Alert {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *AlertCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *AlertCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ac *AlertCreate) check() error {
	if _, ok := ac.mutation.UUID(); !ok {
		return &ValidationError{Name: "UUID", err: errors.New(`ent: missing required field "Alert.UUID"`)}
	}
	if _, ok := ac.mutation.IncidentID(); !ok {
		return &ValidationError{Name: "IncidentID", err: errors.New(`ent: missing required field "Alert.IncidentID"`)}
	}
	if _, ok := ac.mutation.Name(); !ok {
		return &ValidationError{Name: "Name", err: errors.New(`ent: missing required field "Alert.Name"`)}
	}
	if _, ok := ac.mutation.Time(); !ok {
		return &ValidationError{Name: "Time", err: errors.New(`ent: missing required field "Alert.Time"`)}
	}
	if _, ok := ac.mutation.IntLevel(); !ok {
		return &ValidationError{Name: "IntLevel", err: errors.New(`ent: missing required field "Alert.IntLevel"`)}
	}
	if _, ok := ac.mutation.Username(); !ok {
		return &ValidationError{Name: "Username", err: errors.New(`ent: missing required field "Alert.Username"`)}
	}
	if _, ok := ac.mutation.Region(); !ok {
		return &ValidationError{Name: "Region", err: errors.New(`ent: missing required field "Alert.Region"`)}
	}
	if _, ok := ac.mutation.ProbeOS(); !ok {
		return &ValidationError{Name: "ProbeOS", err: errors.New(`ent: missing required field "Alert.ProbeOS"`)}
	}
	if _, ok := ac.mutation.ProbeHost(); !ok {
		return &ValidationError{Name: "ProbeHost", err: errors.New(`ent: missing required field "Alert.ProbeHost"`)}
	}
	return nil
}

func (ac *AlertCreate) sqlSave(ctx context.Context) (*Alert, error) {
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (ac *AlertCreate) createSpec() (*Alert, *sqlgraph.CreateSpec) {
	var (
		_node = &Alert{config: ac.config}
		_spec = &sqlgraph.CreateSpec{
			Table: alert.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: alert.FieldID,
			},
		}
	)
	_spec.OnConflict = ac.conflict
	if value, ok := ac.mutation.UUID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: alert.FieldUUID,
		})
		_node.UUID = value
	}
	if value, ok := ac.mutation.IncidentID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: alert.FieldIncidentID,
		})
		_node.IncidentID = value
	}
	if value, ok := ac.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldName,
		})
		_node.Name = value
	}
	if value, ok := ac.mutation.Time(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: alert.FieldTime,
		})
		_node.Time = value
	}
	if value, ok := ac.mutation.IntLevel(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: alert.FieldIntLevel,
		})
		_node.IntLevel = value
	}
	if value, ok := ac.mutation.Username(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldUsername,
		})
		_node.Username = value
	}
	if value, ok := ac.mutation.Region(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldRegion,
		})
		_node.Region = value
	}
	if value, ok := ac.mutation.ProbeOS(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldProbeOS,
		})
		_node.ProbeOS = value
	}
	if value, ok := ac.mutation.ProbeHost(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldProbeHost,
		})
		_node.ProbeHost = value
	}
	if value, ok := ac.mutation.Error(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldError,
		})
		_node.Error = value
	}
	if nodes := ac.mutation.CountersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   alert.CountersTable,
			Columns: []string{alert.CountersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: counter.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.StatiIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   alert.StatiTable,
			Columns: []string{alert.StatiColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: status.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.FailuresIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   alert.FailuresTable,
			Columns: []string{alert.FailuresColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: failure.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   alert.FilesTable,
			Columns: []string{alert.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Alert.Create().
//		SetUUID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AlertUpsert) {
//			SetUUID(v+v).
//		}).
//		Exec(ctx)
func (ac *AlertCreate) OnConflict(opts ...sql.ConflictOption) *AlertUpsertOne {
	ac.conflict = opts
	return &AlertUpsertOne{
		create: ac,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Alert.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ac *AlertCreate) OnConflictColumns(columns ...string) *AlertUpsertOne {
	ac.conflict = append(ac.conflict, sql.ConflictColumns(columns...))
	return &AlertUpsertOne{
		create: ac,
	}
}

type (
	// AlertUpsertOne is the builder for "upsert"-ing
	//  one Alert node.
	AlertUpsertOne struct {
		create *AlertCreate
	}

	// AlertUpsert is the "OnConflict" setter.
	AlertUpsert struct {
		*sql.UpdateSet
	}
)

// SetUUID sets the "UUID" field.
func (u *AlertUpsert) SetUUID(v uuid.UUID) *AlertUpsert {
	u.Set(alert.FieldUUID, v)
	return u
}

// UpdateUUID sets the "UUID" field to the value that was provided on create.
func (u *AlertUpsert) UpdateUUID() *AlertUpsert {
	u.SetExcluded(alert.FieldUUID)
	return u
}

// SetIncidentID sets the "IncidentID" field.
func (u *AlertUpsert) SetIncidentID(v uuid.UUID) *AlertUpsert {
	u.Set(alert.FieldIncidentID, v)
	return u
}

// UpdateIncidentID sets the "IncidentID" field to the value that was provided on create.
func (u *AlertUpsert) UpdateIncidentID() *AlertUpsert {
	u.SetExcluded(alert.FieldIncidentID)
	return u
}

// SetName sets the "Name" field.
func (u *AlertUpsert) SetName(v string) *AlertUpsert {
	u.Set(alert.FieldName, v)
	return u
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *AlertUpsert) UpdateName() *AlertUpsert {
	u.SetExcluded(alert.FieldName)
	return u
}

// SetTime sets the "Time" field.
func (u *AlertUpsert) SetTime(v time.Time) *AlertUpsert {
	u.Set(alert.FieldTime, v)
	return u
}

// UpdateTime sets the "Time" field to the value that was provided on create.
func (u *AlertUpsert) UpdateTime() *AlertUpsert {
	u.SetExcluded(alert.FieldTime)
	return u
}

// SetIntLevel sets the "IntLevel" field.
func (u *AlertUpsert) SetIntLevel(v int) *AlertUpsert {
	u.Set(alert.FieldIntLevel, v)
	return u
}

// UpdateIntLevel sets the "IntLevel" field to the value that was provided on create.
func (u *AlertUpsert) UpdateIntLevel() *AlertUpsert {
	u.SetExcluded(alert.FieldIntLevel)
	return u
}

// AddIntLevel adds v to the "IntLevel" field.
func (u *AlertUpsert) AddIntLevel(v int) *AlertUpsert {
	u.Add(alert.FieldIntLevel, v)
	return u
}

// SetUsername sets the "Username" field.
func (u *AlertUpsert) SetUsername(v string) *AlertUpsert {
	u.Set(alert.FieldUsername, v)
	return u
}

// UpdateUsername sets the "Username" field to the value that was provided on create.
func (u *AlertUpsert) UpdateUsername() *AlertUpsert {
	u.SetExcluded(alert.FieldUsername)
	return u
}

// SetRegion sets the "Region" field.
func (u *AlertUpsert) SetRegion(v string) *AlertUpsert {
	u.Set(alert.FieldRegion, v)
	return u
}

// UpdateRegion sets the "Region" field to the value that was provided on create.
func (u *AlertUpsert) UpdateRegion() *AlertUpsert {
	u.SetExcluded(alert.FieldRegion)
	return u
}

// SetProbeOS sets the "ProbeOS" field.
func (u *AlertUpsert) SetProbeOS(v string) *AlertUpsert {
	u.Set(alert.FieldProbeOS, v)
	return u
}

// UpdateProbeOS sets the "ProbeOS" field to the value that was provided on create.
func (u *AlertUpsert) UpdateProbeOS() *AlertUpsert {
	u.SetExcluded(alert.FieldProbeOS)
	return u
}

// SetProbeHost sets the "ProbeHost" field.
func (u *AlertUpsert) SetProbeHost(v string) *AlertUpsert {
	u.Set(alert.FieldProbeHost, v)
	return u
}

// UpdateProbeHost sets the "ProbeHost" field to the value that was provided on create.
func (u *AlertUpsert) UpdateProbeHost() *AlertUpsert {
	u.SetExcluded(alert.FieldProbeHost)
	return u
}

// SetError sets the "Error" field.
func (u *AlertUpsert) SetError(v string) *AlertUpsert {
	u.Set(alert.FieldError, v)
	return u
}

// UpdateError sets the "Error" field to the value that was provided on create.
func (u *AlertUpsert) UpdateError() *AlertUpsert {
	u.SetExcluded(alert.FieldError)
	return u
}

// ClearError clears the value of the "Error" field.
func (u *AlertUpsert) ClearError() *AlertUpsert {
	u.SetNull(alert.FieldError)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Alert.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AlertUpsertOne) UpdateNewValues() *AlertUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Alert.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *AlertUpsertOne) Ignore() *AlertUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AlertUpsertOne) DoNothing() *AlertUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AlertCreate.OnConflict
// documentation for more info.
func (u *AlertUpsertOne) Update(set func(*AlertUpsert)) *AlertUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AlertUpsert{UpdateSet: update})
	}))
	return u
}

// SetUUID sets the "UUID" field.
func (u *AlertUpsertOne) SetUUID(v uuid.UUID) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetUUID(v)
	})
}

// UpdateUUID sets the "UUID" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateUUID() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateUUID()
	})
}

// SetIncidentID sets the "IncidentID" field.
func (u *AlertUpsertOne) SetIncidentID(v uuid.UUID) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetIncidentID(v)
	})
}

// UpdateIncidentID sets the "IncidentID" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateIncidentID() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateIncidentID()
	})
}

// SetName sets the "Name" field.
func (u *AlertUpsertOne) SetName(v string) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateName() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateName()
	})
}

// SetTime sets the "Time" field.
func (u *AlertUpsertOne) SetTime(v time.Time) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetTime(v)
	})
}

// UpdateTime sets the "Time" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateTime() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateTime()
	})
}

// SetIntLevel sets the "IntLevel" field.
func (u *AlertUpsertOne) SetIntLevel(v int) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetIntLevel(v)
	})
}

// AddIntLevel adds v to the "IntLevel" field.
func (u *AlertUpsertOne) AddIntLevel(v int) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.AddIntLevel(v)
	})
}

// UpdateIntLevel sets the "IntLevel" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateIntLevel() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateIntLevel()
	})
}

// SetUsername sets the "Username" field.
func (u *AlertUpsertOne) SetUsername(v string) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetUsername(v)
	})
}

// UpdateUsername sets the "Username" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateUsername() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateUsername()
	})
}

// SetRegion sets the "Region" field.
func (u *AlertUpsertOne) SetRegion(v string) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetRegion(v)
	})
}

// UpdateRegion sets the "Region" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateRegion() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateRegion()
	})
}

// SetProbeOS sets the "ProbeOS" field.
func (u *AlertUpsertOne) SetProbeOS(v string) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetProbeOS(v)
	})
}

// UpdateProbeOS sets the "ProbeOS" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateProbeOS() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateProbeOS()
	})
}

// SetProbeHost sets the "ProbeHost" field.
func (u *AlertUpsertOne) SetProbeHost(v string) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetProbeHost(v)
	})
}

// UpdateProbeHost sets the "ProbeHost" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateProbeHost() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateProbeHost()
	})
}

// SetError sets the "Error" field.
func (u *AlertUpsertOne) SetError(v string) *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.SetError(v)
	})
}

// UpdateError sets the "Error" field to the value that was provided on create.
func (u *AlertUpsertOne) UpdateError() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateError()
	})
}

// ClearError clears the value of the "Error" field.
func (u *AlertUpsertOne) ClearError() *AlertUpsertOne {
	return u.Update(func(s *AlertUpsert) {
		s.ClearError()
	})
}

// Exec executes the query.
func (u *AlertUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AlertCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AlertUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *AlertUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *AlertUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// AlertCreateBulk is the builder for creating many Alert entities in bulk.
type AlertCreateBulk struct {
	config
	builders []*AlertCreate
	conflict []sql.ConflictOption
}

// Save creates the Alert entities in the database.
func (acb *AlertCreateBulk) Save(ctx context.Context) ([]*Alert, error) {
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Alert, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AlertMutation)
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
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = acb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *AlertCreateBulk) SaveX(ctx context.Context) []*Alert {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *AlertCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *AlertCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Alert.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AlertUpsert) {
//			SetUUID(v+v).
//		}).
//		Exec(ctx)
func (acb *AlertCreateBulk) OnConflict(opts ...sql.ConflictOption) *AlertUpsertBulk {
	acb.conflict = opts
	return &AlertUpsertBulk{
		create: acb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Alert.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (acb *AlertCreateBulk) OnConflictColumns(columns ...string) *AlertUpsertBulk {
	acb.conflict = append(acb.conflict, sql.ConflictColumns(columns...))
	return &AlertUpsertBulk{
		create: acb,
	}
}

// AlertUpsertBulk is the builder for "upsert"-ing
// a bulk of Alert nodes.
type AlertUpsertBulk struct {
	create *AlertCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Alert.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AlertUpsertBulk) UpdateNewValues() *AlertUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Alert.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *AlertUpsertBulk) Ignore() *AlertUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AlertUpsertBulk) DoNothing() *AlertUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AlertCreateBulk.OnConflict
// documentation for more info.
func (u *AlertUpsertBulk) Update(set func(*AlertUpsert)) *AlertUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AlertUpsert{UpdateSet: update})
	}))
	return u
}

// SetUUID sets the "UUID" field.
func (u *AlertUpsertBulk) SetUUID(v uuid.UUID) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetUUID(v)
	})
}

// UpdateUUID sets the "UUID" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateUUID() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateUUID()
	})
}

// SetIncidentID sets the "IncidentID" field.
func (u *AlertUpsertBulk) SetIncidentID(v uuid.UUID) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetIncidentID(v)
	})
}

// UpdateIncidentID sets the "IncidentID" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateIncidentID() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateIncidentID()
	})
}

// SetName sets the "Name" field.
func (u *AlertUpsertBulk) SetName(v string) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "Name" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateName() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateName()
	})
}

// SetTime sets the "Time" field.
func (u *AlertUpsertBulk) SetTime(v time.Time) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetTime(v)
	})
}

// UpdateTime sets the "Time" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateTime() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateTime()
	})
}

// SetIntLevel sets the "IntLevel" field.
func (u *AlertUpsertBulk) SetIntLevel(v int) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetIntLevel(v)
	})
}

// AddIntLevel adds v to the "IntLevel" field.
func (u *AlertUpsertBulk) AddIntLevel(v int) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.AddIntLevel(v)
	})
}

// UpdateIntLevel sets the "IntLevel" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateIntLevel() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateIntLevel()
	})
}

// SetUsername sets the "Username" field.
func (u *AlertUpsertBulk) SetUsername(v string) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetUsername(v)
	})
}

// UpdateUsername sets the "Username" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateUsername() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateUsername()
	})
}

// SetRegion sets the "Region" field.
func (u *AlertUpsertBulk) SetRegion(v string) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetRegion(v)
	})
}

// UpdateRegion sets the "Region" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateRegion() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateRegion()
	})
}

// SetProbeOS sets the "ProbeOS" field.
func (u *AlertUpsertBulk) SetProbeOS(v string) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetProbeOS(v)
	})
}

// UpdateProbeOS sets the "ProbeOS" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateProbeOS() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateProbeOS()
	})
}

// SetProbeHost sets the "ProbeHost" field.
func (u *AlertUpsertBulk) SetProbeHost(v string) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetProbeHost(v)
	})
}

// UpdateProbeHost sets the "ProbeHost" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateProbeHost() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateProbeHost()
	})
}

// SetError sets the "Error" field.
func (u *AlertUpsertBulk) SetError(v string) *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.SetError(v)
	})
}

// UpdateError sets the "Error" field to the value that was provided on create.
func (u *AlertUpsertBulk) UpdateError() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.UpdateError()
	})
}

// ClearError clears the value of the "Error" field.
func (u *AlertUpsertBulk) ClearError() *AlertUpsertBulk {
	return u.Update(func(s *AlertUpsert) {
		s.ClearError()
	})
}

// Exec executes the query.
func (u *AlertUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the AlertCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AlertCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AlertUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
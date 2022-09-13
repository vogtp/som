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
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/counter"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/status"
)

// IncidentUpdate is the builder for updating Incident entities.
type IncidentUpdate struct {
	config
	hooks    []Hook
	mutation *IncidentMutation
}

// Where appends a list predicates to the IncidentUpdate builder.
func (iu *IncidentUpdate) Where(ps ...predicate.Incident) *IncidentUpdate {
	iu.mutation.Where(ps...)
	return iu
}

// SetUUID sets the "UUID" field.
func (iu *IncidentUpdate) SetUUID(u uuid.UUID) *IncidentUpdate {
	iu.mutation.SetUUID(u)
	return iu
}

// SetIncidentID sets the "IncidentID" field.
func (iu *IncidentUpdate) SetIncidentID(u uuid.UUID) *IncidentUpdate {
	iu.mutation.SetIncidentID(u)
	return iu
}

// SetName sets the "Name" field.
func (iu *IncidentUpdate) SetName(s string) *IncidentUpdate {
	iu.mutation.SetName(s)
	return iu
}

// SetTime sets the "Time" field.
func (iu *IncidentUpdate) SetTime(t time.Time) *IncidentUpdate {
	iu.mutation.SetTime(t)
	return iu
}

// SetIntLevel sets the "IntLevel" field.
func (iu *IncidentUpdate) SetIntLevel(i int) *IncidentUpdate {
	iu.mutation.ResetIntLevel()
	iu.mutation.SetIntLevel(i)
	return iu
}

// AddIntLevel adds i to the "IntLevel" field.
func (iu *IncidentUpdate) AddIntLevel(i int) *IncidentUpdate {
	iu.mutation.AddIntLevel(i)
	return iu
}

// SetUsername sets the "Username" field.
func (iu *IncidentUpdate) SetUsername(s string) *IncidentUpdate {
	iu.mutation.SetUsername(s)
	return iu
}

// SetRegion sets the "Region" field.
func (iu *IncidentUpdate) SetRegion(s string) *IncidentUpdate {
	iu.mutation.SetRegion(s)
	return iu
}

// SetProbeOS sets the "ProbeOS" field.
func (iu *IncidentUpdate) SetProbeOS(s string) *IncidentUpdate {
	iu.mutation.SetProbeOS(s)
	return iu
}

// SetProbeHost sets the "ProbeHost" field.
func (iu *IncidentUpdate) SetProbeHost(s string) *IncidentUpdate {
	iu.mutation.SetProbeHost(s)
	return iu
}

// SetError sets the "Error" field.
func (iu *IncidentUpdate) SetError(s string) *IncidentUpdate {
	iu.mutation.SetError(s)
	return iu
}

// SetNillableError sets the "Error" field if the given value is not nil.
func (iu *IncidentUpdate) SetNillableError(s *string) *IncidentUpdate {
	if s != nil {
		iu.SetError(*s)
	}
	return iu
}

// ClearError clears the value of the "Error" field.
func (iu *IncidentUpdate) ClearError() *IncidentUpdate {
	iu.mutation.ClearError()
	return iu
}

// SetStart sets the "Start" field.
func (iu *IncidentUpdate) SetStart(t time.Time) *IncidentUpdate {
	iu.mutation.SetStart(t)
	return iu
}

// SetEnd sets the "End" field.
func (iu *IncidentUpdate) SetEnd(t time.Time) *IncidentUpdate {
	iu.mutation.SetEnd(t)
	return iu
}

// SetState sets the "State" field.
func (iu *IncidentUpdate) SetState(b []byte) *IncidentUpdate {
	iu.mutation.SetState(b)
	return iu
}

// AddCounterIDs adds the "Counters" edge to the Counter entity by IDs.
func (iu *IncidentUpdate) AddCounterIDs(ids ...int) *IncidentUpdate {
	iu.mutation.AddCounterIDs(ids...)
	return iu
}

// AddCounters adds the "Counters" edges to the Counter entity.
func (iu *IncidentUpdate) AddCounters(c ...*Counter) *IncidentUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iu.AddCounterIDs(ids...)
}

// AddStatiIDs adds the "Stati" edge to the Status entity by IDs.
func (iu *IncidentUpdate) AddStatiIDs(ids ...int) *IncidentUpdate {
	iu.mutation.AddStatiIDs(ids...)
	return iu
}

// AddStati adds the "Stati" edges to the Status entity.
func (iu *IncidentUpdate) AddStati(s ...*Status) *IncidentUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return iu.AddStatiIDs(ids...)
}

// AddFailureIDs adds the "Failures" edge to the Failure entity by IDs.
func (iu *IncidentUpdate) AddFailureIDs(ids ...int) *IncidentUpdate {
	iu.mutation.AddFailureIDs(ids...)
	return iu
}

// AddFailures adds the "Failures" edges to the Failure entity.
func (iu *IncidentUpdate) AddFailures(f ...*Failure) *IncidentUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iu.AddFailureIDs(ids...)
}

// AddFileIDs adds the "Files" edge to the File entity by IDs.
func (iu *IncidentUpdate) AddFileIDs(ids ...int) *IncidentUpdate {
	iu.mutation.AddFileIDs(ids...)
	return iu
}

// AddFiles adds the "Files" edges to the File entity.
func (iu *IncidentUpdate) AddFiles(f ...*File) *IncidentUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iu.AddFileIDs(ids...)
}

// Mutation returns the IncidentMutation object of the builder.
func (iu *IncidentUpdate) Mutation() *IncidentMutation {
	return iu.mutation
}

// ClearCounters clears all "Counters" edges to the Counter entity.
func (iu *IncidentUpdate) ClearCounters() *IncidentUpdate {
	iu.mutation.ClearCounters()
	return iu
}

// RemoveCounterIDs removes the "Counters" edge to Counter entities by IDs.
func (iu *IncidentUpdate) RemoveCounterIDs(ids ...int) *IncidentUpdate {
	iu.mutation.RemoveCounterIDs(ids...)
	return iu
}

// RemoveCounters removes "Counters" edges to Counter entities.
func (iu *IncidentUpdate) RemoveCounters(c ...*Counter) *IncidentUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iu.RemoveCounterIDs(ids...)
}

// ClearStati clears all "Stati" edges to the Status entity.
func (iu *IncidentUpdate) ClearStati() *IncidentUpdate {
	iu.mutation.ClearStati()
	return iu
}

// RemoveStatiIDs removes the "Stati" edge to Status entities by IDs.
func (iu *IncidentUpdate) RemoveStatiIDs(ids ...int) *IncidentUpdate {
	iu.mutation.RemoveStatiIDs(ids...)
	return iu
}

// RemoveStati removes "Stati" edges to Status entities.
func (iu *IncidentUpdate) RemoveStati(s ...*Status) *IncidentUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return iu.RemoveStatiIDs(ids...)
}

// ClearFailures clears all "Failures" edges to the Failure entity.
func (iu *IncidentUpdate) ClearFailures() *IncidentUpdate {
	iu.mutation.ClearFailures()
	return iu
}

// RemoveFailureIDs removes the "Failures" edge to Failure entities by IDs.
func (iu *IncidentUpdate) RemoveFailureIDs(ids ...int) *IncidentUpdate {
	iu.mutation.RemoveFailureIDs(ids...)
	return iu
}

// RemoveFailures removes "Failures" edges to Failure entities.
func (iu *IncidentUpdate) RemoveFailures(f ...*Failure) *IncidentUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iu.RemoveFailureIDs(ids...)
}

// ClearFiles clears all "Files" edges to the File entity.
func (iu *IncidentUpdate) ClearFiles() *IncidentUpdate {
	iu.mutation.ClearFiles()
	return iu
}

// RemoveFileIDs removes the "Files" edge to File entities by IDs.
func (iu *IncidentUpdate) RemoveFileIDs(ids ...int) *IncidentUpdate {
	iu.mutation.RemoveFileIDs(ids...)
	return iu
}

// RemoveFiles removes "Files" edges to File entities.
func (iu *IncidentUpdate) RemoveFiles(f ...*File) *IncidentUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iu.RemoveFileIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (iu *IncidentUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(iu.hooks) == 0 {
		affected, err = iu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*IncidentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			iu.mutation = mutation
			affected, err = iu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(iu.hooks) - 1; i >= 0; i-- {
			if iu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = iu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, iu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (iu *IncidentUpdate) SaveX(ctx context.Context) int {
	affected, err := iu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (iu *IncidentUpdate) Exec(ctx context.Context) error {
	_, err := iu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iu *IncidentUpdate) ExecX(ctx context.Context) {
	if err := iu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (iu *IncidentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   incident.Table,
			Columns: incident.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: incident.FieldID,
			},
		},
	}
	if ps := iu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iu.mutation.UUID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: incident.FieldUUID,
		})
	}
	if value, ok := iu.mutation.IncidentID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: incident.FieldIncidentID,
		})
	}
	if value, ok := iu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldName,
		})
	}
	if value, ok := iu.mutation.Time(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: incident.FieldTime,
		})
	}
	if value, ok := iu.mutation.IntLevel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: incident.FieldIntLevel,
		})
	}
	if value, ok := iu.mutation.AddedIntLevel(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: incident.FieldIntLevel,
		})
	}
	if value, ok := iu.mutation.Username(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldUsername,
		})
	}
	if value, ok := iu.mutation.Region(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldRegion,
		})
	}
	if value, ok := iu.mutation.ProbeOS(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldProbeOS,
		})
	}
	if value, ok := iu.mutation.ProbeHost(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldProbeHost,
		})
	}
	if value, ok := iu.mutation.Error(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldError,
		})
	}
	if iu.mutation.ErrorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: incident.FieldError,
		})
	}
	if value, ok := iu.mutation.Start(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: incident.FieldStart,
		})
	}
	if value, ok := iu.mutation.End(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: incident.FieldEnd,
		})
	}
	if value, ok := iu.mutation.State(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBytes,
			Value:  value,
			Column: incident.FieldState,
		})
	}
	if iu.mutation.CountersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.CountersTable,
			Columns: []string{incident.CountersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: counter.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedCountersIDs(); len(nodes) > 0 && !iu.mutation.CountersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.CountersTable,
			Columns: []string{incident.CountersColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.CountersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.CountersTable,
			Columns: []string{incident.CountersColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iu.mutation.StatiCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.StatiTable,
			Columns: []string{incident.StatiColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: status.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedStatiIDs(); len(nodes) > 0 && !iu.mutation.StatiCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.StatiTable,
			Columns: []string{incident.StatiColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.StatiIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.StatiTable,
			Columns: []string{incident.StatiColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iu.mutation.FailuresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FailuresTable,
			Columns: []string{incident.FailuresColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: failure.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedFailuresIDs(); len(nodes) > 0 && !iu.mutation.FailuresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FailuresTable,
			Columns: []string{incident.FailuresColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.FailuresIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FailuresTable,
			Columns: []string{incident.FailuresColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iu.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FilesTable,
			Columns: []string{incident.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedFilesIDs(); len(nodes) > 0 && !iu.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FilesTable,
			Columns: []string{incident.FilesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FilesTable,
			Columns: []string{incident.FilesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, iu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{incident.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// IncidentUpdateOne is the builder for updating a single Incident entity.
type IncidentUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *IncidentMutation
}

// SetUUID sets the "UUID" field.
func (iuo *IncidentUpdateOne) SetUUID(u uuid.UUID) *IncidentUpdateOne {
	iuo.mutation.SetUUID(u)
	return iuo
}

// SetIncidentID sets the "IncidentID" field.
func (iuo *IncidentUpdateOne) SetIncidentID(u uuid.UUID) *IncidentUpdateOne {
	iuo.mutation.SetIncidentID(u)
	return iuo
}

// SetName sets the "Name" field.
func (iuo *IncidentUpdateOne) SetName(s string) *IncidentUpdateOne {
	iuo.mutation.SetName(s)
	return iuo
}

// SetTime sets the "Time" field.
func (iuo *IncidentUpdateOne) SetTime(t time.Time) *IncidentUpdateOne {
	iuo.mutation.SetTime(t)
	return iuo
}

// SetIntLevel sets the "IntLevel" field.
func (iuo *IncidentUpdateOne) SetIntLevel(i int) *IncidentUpdateOne {
	iuo.mutation.ResetIntLevel()
	iuo.mutation.SetIntLevel(i)
	return iuo
}

// AddIntLevel adds i to the "IntLevel" field.
func (iuo *IncidentUpdateOne) AddIntLevel(i int) *IncidentUpdateOne {
	iuo.mutation.AddIntLevel(i)
	return iuo
}

// SetUsername sets the "Username" field.
func (iuo *IncidentUpdateOne) SetUsername(s string) *IncidentUpdateOne {
	iuo.mutation.SetUsername(s)
	return iuo
}

// SetRegion sets the "Region" field.
func (iuo *IncidentUpdateOne) SetRegion(s string) *IncidentUpdateOne {
	iuo.mutation.SetRegion(s)
	return iuo
}

// SetProbeOS sets the "ProbeOS" field.
func (iuo *IncidentUpdateOne) SetProbeOS(s string) *IncidentUpdateOne {
	iuo.mutation.SetProbeOS(s)
	return iuo
}

// SetProbeHost sets the "ProbeHost" field.
func (iuo *IncidentUpdateOne) SetProbeHost(s string) *IncidentUpdateOne {
	iuo.mutation.SetProbeHost(s)
	return iuo
}

// SetError sets the "Error" field.
func (iuo *IncidentUpdateOne) SetError(s string) *IncidentUpdateOne {
	iuo.mutation.SetError(s)
	return iuo
}

// SetNillableError sets the "Error" field if the given value is not nil.
func (iuo *IncidentUpdateOne) SetNillableError(s *string) *IncidentUpdateOne {
	if s != nil {
		iuo.SetError(*s)
	}
	return iuo
}

// ClearError clears the value of the "Error" field.
func (iuo *IncidentUpdateOne) ClearError() *IncidentUpdateOne {
	iuo.mutation.ClearError()
	return iuo
}

// SetStart sets the "Start" field.
func (iuo *IncidentUpdateOne) SetStart(t time.Time) *IncidentUpdateOne {
	iuo.mutation.SetStart(t)
	return iuo
}

// SetEnd sets the "End" field.
func (iuo *IncidentUpdateOne) SetEnd(t time.Time) *IncidentUpdateOne {
	iuo.mutation.SetEnd(t)
	return iuo
}

// SetState sets the "State" field.
func (iuo *IncidentUpdateOne) SetState(b []byte) *IncidentUpdateOne {
	iuo.mutation.SetState(b)
	return iuo
}

// AddCounterIDs adds the "Counters" edge to the Counter entity by IDs.
func (iuo *IncidentUpdateOne) AddCounterIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.AddCounterIDs(ids...)
	return iuo
}

// AddCounters adds the "Counters" edges to the Counter entity.
func (iuo *IncidentUpdateOne) AddCounters(c ...*Counter) *IncidentUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iuo.AddCounterIDs(ids...)
}

// AddStatiIDs adds the "Stati" edge to the Status entity by IDs.
func (iuo *IncidentUpdateOne) AddStatiIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.AddStatiIDs(ids...)
	return iuo
}

// AddStati adds the "Stati" edges to the Status entity.
func (iuo *IncidentUpdateOne) AddStati(s ...*Status) *IncidentUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return iuo.AddStatiIDs(ids...)
}

// AddFailureIDs adds the "Failures" edge to the Failure entity by IDs.
func (iuo *IncidentUpdateOne) AddFailureIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.AddFailureIDs(ids...)
	return iuo
}

// AddFailures adds the "Failures" edges to the Failure entity.
func (iuo *IncidentUpdateOne) AddFailures(f ...*Failure) *IncidentUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iuo.AddFailureIDs(ids...)
}

// AddFileIDs adds the "Files" edge to the File entity by IDs.
func (iuo *IncidentUpdateOne) AddFileIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.AddFileIDs(ids...)
	return iuo
}

// AddFiles adds the "Files" edges to the File entity.
func (iuo *IncidentUpdateOne) AddFiles(f ...*File) *IncidentUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iuo.AddFileIDs(ids...)
}

// Mutation returns the IncidentMutation object of the builder.
func (iuo *IncidentUpdateOne) Mutation() *IncidentMutation {
	return iuo.mutation
}

// ClearCounters clears all "Counters" edges to the Counter entity.
func (iuo *IncidentUpdateOne) ClearCounters() *IncidentUpdateOne {
	iuo.mutation.ClearCounters()
	return iuo
}

// RemoveCounterIDs removes the "Counters" edge to Counter entities by IDs.
func (iuo *IncidentUpdateOne) RemoveCounterIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.RemoveCounterIDs(ids...)
	return iuo
}

// RemoveCounters removes "Counters" edges to Counter entities.
func (iuo *IncidentUpdateOne) RemoveCounters(c ...*Counter) *IncidentUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iuo.RemoveCounterIDs(ids...)
}

// ClearStati clears all "Stati" edges to the Status entity.
func (iuo *IncidentUpdateOne) ClearStati() *IncidentUpdateOne {
	iuo.mutation.ClearStati()
	return iuo
}

// RemoveStatiIDs removes the "Stati" edge to Status entities by IDs.
func (iuo *IncidentUpdateOne) RemoveStatiIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.RemoveStatiIDs(ids...)
	return iuo
}

// RemoveStati removes "Stati" edges to Status entities.
func (iuo *IncidentUpdateOne) RemoveStati(s ...*Status) *IncidentUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return iuo.RemoveStatiIDs(ids...)
}

// ClearFailures clears all "Failures" edges to the Failure entity.
func (iuo *IncidentUpdateOne) ClearFailures() *IncidentUpdateOne {
	iuo.mutation.ClearFailures()
	return iuo
}

// RemoveFailureIDs removes the "Failures" edge to Failure entities by IDs.
func (iuo *IncidentUpdateOne) RemoveFailureIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.RemoveFailureIDs(ids...)
	return iuo
}

// RemoveFailures removes "Failures" edges to Failure entities.
func (iuo *IncidentUpdateOne) RemoveFailures(f ...*Failure) *IncidentUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iuo.RemoveFailureIDs(ids...)
}

// ClearFiles clears all "Files" edges to the File entity.
func (iuo *IncidentUpdateOne) ClearFiles() *IncidentUpdateOne {
	iuo.mutation.ClearFiles()
	return iuo
}

// RemoveFileIDs removes the "Files" edge to File entities by IDs.
func (iuo *IncidentUpdateOne) RemoveFileIDs(ids ...int) *IncidentUpdateOne {
	iuo.mutation.RemoveFileIDs(ids...)
	return iuo
}

// RemoveFiles removes "Files" edges to File entities.
func (iuo *IncidentUpdateOne) RemoveFiles(f ...*File) *IncidentUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return iuo.RemoveFileIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (iuo *IncidentUpdateOne) Select(field string, fields ...string) *IncidentUpdateOne {
	iuo.fields = append([]string{field}, fields...)
	return iuo
}

// Save executes the query and returns the updated Incident entity.
func (iuo *IncidentUpdateOne) Save(ctx context.Context) (*Incident, error) {
	var (
		err  error
		node *Incident
	)
	if len(iuo.hooks) == 0 {
		node, err = iuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*IncidentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			iuo.mutation = mutation
			node, err = iuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(iuo.hooks) - 1; i >= 0; i-- {
			if iuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = iuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, iuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Incident)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from IncidentMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (iuo *IncidentUpdateOne) SaveX(ctx context.Context) *Incident {
	node, err := iuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (iuo *IncidentUpdateOne) Exec(ctx context.Context) error {
	_, err := iuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iuo *IncidentUpdateOne) ExecX(ctx context.Context) {
	if err := iuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (iuo *IncidentUpdateOne) sqlSave(ctx context.Context) (_node *Incident, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   incident.Table,
			Columns: incident.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: incident.FieldID,
			},
		},
	}
	id, ok := iuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Incident.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := iuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, incident.FieldID)
		for _, f := range fields {
			if !incident.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != incident.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := iuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iuo.mutation.UUID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: incident.FieldUUID,
		})
	}
	if value, ok := iuo.mutation.IncidentID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: incident.FieldIncidentID,
		})
	}
	if value, ok := iuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldName,
		})
	}
	if value, ok := iuo.mutation.Time(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: incident.FieldTime,
		})
	}
	if value, ok := iuo.mutation.IntLevel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: incident.FieldIntLevel,
		})
	}
	if value, ok := iuo.mutation.AddedIntLevel(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: incident.FieldIntLevel,
		})
	}
	if value, ok := iuo.mutation.Username(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldUsername,
		})
	}
	if value, ok := iuo.mutation.Region(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldRegion,
		})
	}
	if value, ok := iuo.mutation.ProbeOS(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldProbeOS,
		})
	}
	if value, ok := iuo.mutation.ProbeHost(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldProbeHost,
		})
	}
	if value, ok := iuo.mutation.Error(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: incident.FieldError,
		})
	}
	if iuo.mutation.ErrorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: incident.FieldError,
		})
	}
	if value, ok := iuo.mutation.Start(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: incident.FieldStart,
		})
	}
	if value, ok := iuo.mutation.End(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: incident.FieldEnd,
		})
	}
	if value, ok := iuo.mutation.State(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBytes,
			Value:  value,
			Column: incident.FieldState,
		})
	}
	if iuo.mutation.CountersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.CountersTable,
			Columns: []string{incident.CountersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: counter.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedCountersIDs(); len(nodes) > 0 && !iuo.mutation.CountersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.CountersTable,
			Columns: []string{incident.CountersColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.CountersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.CountersTable,
			Columns: []string{incident.CountersColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iuo.mutation.StatiCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.StatiTable,
			Columns: []string{incident.StatiColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: status.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedStatiIDs(); len(nodes) > 0 && !iuo.mutation.StatiCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.StatiTable,
			Columns: []string{incident.StatiColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.StatiIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.StatiTable,
			Columns: []string{incident.StatiColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iuo.mutation.FailuresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FailuresTable,
			Columns: []string{incident.FailuresColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: failure.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedFailuresIDs(); len(nodes) > 0 && !iuo.mutation.FailuresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FailuresTable,
			Columns: []string{incident.FailuresColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.FailuresIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FailuresTable,
			Columns: []string{incident.FailuresColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if iuo.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FilesTable,
			Columns: []string{incident.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedFilesIDs(); len(nodes) > 0 && !iuo.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FilesTable,
			Columns: []string{incident.FilesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   incident.FilesTable,
			Columns: []string{incident.FilesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Incident{config: iuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, iuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{incident.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}
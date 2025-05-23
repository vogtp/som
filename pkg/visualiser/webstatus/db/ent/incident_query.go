// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/counter"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/status"
)

// IncidentQuery is the builder for querying Incident entities.
type IncidentQuery struct {
	config
	ctx               *QueryContext
	order             []incident.OrderOption
	inters            []Interceptor
	predicates        []predicate.Incident
	withCounters      *CounterQuery
	withStati         *StatusQuery
	withFailures      *FailureQuery
	withFiles         *FileQuery
	modifiers         []func(*sql.Selector)
	loadTotal         []func(context.Context, []*Incident) error
	withNamedCounters map[string]*CounterQuery
	withNamedStati    map[string]*StatusQuery
	withNamedFailures map[string]*FailureQuery
	withNamedFiles    map[string]*FileQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the IncidentQuery builder.
func (iq *IncidentQuery) Where(ps ...predicate.Incident) *IncidentQuery {
	iq.predicates = append(iq.predicates, ps...)
	return iq
}

// Limit the number of records to be returned by this query.
func (iq *IncidentQuery) Limit(limit int) *IncidentQuery {
	iq.ctx.Limit = &limit
	return iq
}

// Offset to start from.
func (iq *IncidentQuery) Offset(offset int) *IncidentQuery {
	iq.ctx.Offset = &offset
	return iq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (iq *IncidentQuery) Unique(unique bool) *IncidentQuery {
	iq.ctx.Unique = &unique
	return iq
}

// Order specifies how the records should be ordered.
func (iq *IncidentQuery) Order(o ...incident.OrderOption) *IncidentQuery {
	iq.order = append(iq.order, o...)
	return iq
}

// QueryCounters chains the current query on the "Counters" edge.
func (iq *IncidentQuery) QueryCounters() *CounterQuery {
	query := (&CounterClient{config: iq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := iq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := iq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(incident.Table, incident.FieldID, selector),
			sqlgraph.To(counter.Table, counter.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, incident.CountersTable, incident.CountersColumn),
		)
		fromU = sqlgraph.SetNeighbors(iq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryStati chains the current query on the "Stati" edge.
func (iq *IncidentQuery) QueryStati() *StatusQuery {
	query := (&StatusClient{config: iq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := iq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := iq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(incident.Table, incident.FieldID, selector),
			sqlgraph.To(status.Table, status.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, incident.StatiTable, incident.StatiColumn),
		)
		fromU = sqlgraph.SetNeighbors(iq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFailures chains the current query on the "Failures" edge.
func (iq *IncidentQuery) QueryFailures() *FailureQuery {
	query := (&FailureClient{config: iq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := iq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := iq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(incident.Table, incident.FieldID, selector),
			sqlgraph.To(failure.Table, failure.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, incident.FailuresTable, incident.FailuresColumn),
		)
		fromU = sqlgraph.SetNeighbors(iq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFiles chains the current query on the "Files" edge.
func (iq *IncidentQuery) QueryFiles() *FileQuery {
	query := (&FileClient{config: iq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := iq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := iq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(incident.Table, incident.FieldID, selector),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, incident.FilesTable, incident.FilesColumn),
		)
		fromU = sqlgraph.SetNeighbors(iq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Incident entity from the query.
// Returns a *NotFoundError when no Incident was found.
func (iq *IncidentQuery) First(ctx context.Context) (*Incident, error) {
	nodes, err := iq.Limit(1).All(setContextOp(ctx, iq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{incident.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (iq *IncidentQuery) FirstX(ctx context.Context) *Incident {
	node, err := iq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Incident ID from the query.
// Returns a *NotFoundError when no Incident ID was found.
func (iq *IncidentQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = iq.Limit(1).IDs(setContextOp(ctx, iq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{incident.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (iq *IncidentQuery) FirstIDX(ctx context.Context) int {
	id, err := iq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Incident entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Incident entity is found.
// Returns a *NotFoundError when no Incident entities are found.
func (iq *IncidentQuery) Only(ctx context.Context) (*Incident, error) {
	nodes, err := iq.Limit(2).All(setContextOp(ctx, iq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{incident.Label}
	default:
		return nil, &NotSingularError{incident.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (iq *IncidentQuery) OnlyX(ctx context.Context) *Incident {
	node, err := iq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Incident ID in the query.
// Returns a *NotSingularError when more than one Incident ID is found.
// Returns a *NotFoundError when no entities are found.
func (iq *IncidentQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = iq.Limit(2).IDs(setContextOp(ctx, iq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{incident.Label}
	default:
		err = &NotSingularError{incident.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (iq *IncidentQuery) OnlyIDX(ctx context.Context) int {
	id, err := iq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Incidents.
func (iq *IncidentQuery) All(ctx context.Context) ([]*Incident, error) {
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryAll)
	if err := iq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Incident, *IncidentQuery]()
	return withInterceptors[[]*Incident](ctx, iq, qr, iq.inters)
}

// AllX is like All, but panics if an error occurs.
func (iq *IncidentQuery) AllX(ctx context.Context) []*Incident {
	nodes, err := iq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Incident IDs.
func (iq *IncidentQuery) IDs(ctx context.Context) (ids []int, err error) {
	if iq.ctx.Unique == nil && iq.path != nil {
		iq.Unique(true)
	}
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryIDs)
	if err = iq.Select(incident.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (iq *IncidentQuery) IDsX(ctx context.Context) []int {
	ids, err := iq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (iq *IncidentQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryCount)
	if err := iq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, iq, querierCount[*IncidentQuery](), iq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (iq *IncidentQuery) CountX(ctx context.Context) int {
	count, err := iq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (iq *IncidentQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryExist)
	switch _, err := iq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (iq *IncidentQuery) ExistX(ctx context.Context) bool {
	exist, err := iq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the IncidentQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (iq *IncidentQuery) Clone() *IncidentQuery {
	if iq == nil {
		return nil
	}
	return &IncidentQuery{
		config:       iq.config,
		ctx:          iq.ctx.Clone(),
		order:        append([]incident.OrderOption{}, iq.order...),
		inters:       append([]Interceptor{}, iq.inters...),
		predicates:   append([]predicate.Incident{}, iq.predicates...),
		withCounters: iq.withCounters.Clone(),
		withStati:    iq.withStati.Clone(),
		withFailures: iq.withFailures.Clone(),
		withFiles:    iq.withFiles.Clone(),
		// clone intermediate query.
		sql:  iq.sql.Clone(),
		path: iq.path,
	}
}

// WithCounters tells the query-builder to eager-load the nodes that are connected to
// the "Counters" edge. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithCounters(opts ...func(*CounterQuery)) *IncidentQuery {
	query := (&CounterClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	iq.withCounters = query
	return iq
}

// WithStati tells the query-builder to eager-load the nodes that are connected to
// the "Stati" edge. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithStati(opts ...func(*StatusQuery)) *IncidentQuery {
	query := (&StatusClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	iq.withStati = query
	return iq
}

// WithFailures tells the query-builder to eager-load the nodes that are connected to
// the "Failures" edge. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithFailures(opts ...func(*FailureQuery)) *IncidentQuery {
	query := (&FailureClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	iq.withFailures = query
	return iq
}

// WithFiles tells the query-builder to eager-load the nodes that are connected to
// the "Files" edge. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithFiles(opts ...func(*FileQuery)) *IncidentQuery {
	query := (&FileClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	iq.withFiles = query
	return iq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		UUID uuid.UUID `json:"UUID,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Incident.Query().
//		GroupBy(incident.FieldUUID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (iq *IncidentQuery) GroupBy(field string, fields ...string) *IncidentGroupBy {
	iq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &IncidentGroupBy{build: iq}
	grbuild.flds = &iq.ctx.Fields
	grbuild.label = incident.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		UUID uuid.UUID `json:"UUID,omitempty"`
//	}
//
//	client.Incident.Query().
//		Select(incident.FieldUUID).
//		Scan(ctx, &v)
func (iq *IncidentQuery) Select(fields ...string) *IncidentSelect {
	iq.ctx.Fields = append(iq.ctx.Fields, fields...)
	sbuild := &IncidentSelect{IncidentQuery: iq}
	sbuild.label = incident.Label
	sbuild.flds, sbuild.scan = &iq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a IncidentSelect configured with the given aggregations.
func (iq *IncidentQuery) Aggregate(fns ...AggregateFunc) *IncidentSelect {
	return iq.Select().Aggregate(fns...)
}

func (iq *IncidentQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range iq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, iq); err != nil {
				return err
			}
		}
	}
	for _, f := range iq.ctx.Fields {
		if !incident.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if iq.path != nil {
		prev, err := iq.path(ctx)
		if err != nil {
			return err
		}
		iq.sql = prev
	}
	return nil
}

func (iq *IncidentQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Incident, error) {
	var (
		nodes       = []*Incident{}
		_spec       = iq.querySpec()
		loadedTypes = [4]bool{
			iq.withCounters != nil,
			iq.withStati != nil,
			iq.withFailures != nil,
			iq.withFiles != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Incident).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Incident{config: iq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(iq.modifiers) > 0 {
		_spec.Modifiers = iq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, iq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := iq.withCounters; query != nil {
		if err := iq.loadCounters(ctx, query, nodes,
			func(n *Incident) { n.Edges.Counters = []*Counter{} },
			func(n *Incident, e *Counter) { n.Edges.Counters = append(n.Edges.Counters, e) }); err != nil {
			return nil, err
		}
	}
	if query := iq.withStati; query != nil {
		if err := iq.loadStati(ctx, query, nodes,
			func(n *Incident) { n.Edges.Stati = []*Status{} },
			func(n *Incident, e *Status) { n.Edges.Stati = append(n.Edges.Stati, e) }); err != nil {
			return nil, err
		}
	}
	if query := iq.withFailures; query != nil {
		if err := iq.loadFailures(ctx, query, nodes,
			func(n *Incident) { n.Edges.Failures = []*Failure{} },
			func(n *Incident, e *Failure) { n.Edges.Failures = append(n.Edges.Failures, e) }); err != nil {
			return nil, err
		}
	}
	if query := iq.withFiles; query != nil {
		if err := iq.loadFiles(ctx, query, nodes,
			func(n *Incident) { n.Edges.Files = []*File{} },
			func(n *Incident, e *File) { n.Edges.Files = append(n.Edges.Files, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range iq.withNamedCounters {
		if err := iq.loadCounters(ctx, query, nodes,
			func(n *Incident) { n.appendNamedCounters(name) },
			func(n *Incident, e *Counter) { n.appendNamedCounters(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range iq.withNamedStati {
		if err := iq.loadStati(ctx, query, nodes,
			func(n *Incident) { n.appendNamedStati(name) },
			func(n *Incident, e *Status) { n.appendNamedStati(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range iq.withNamedFailures {
		if err := iq.loadFailures(ctx, query, nodes,
			func(n *Incident) { n.appendNamedFailures(name) },
			func(n *Incident, e *Failure) { n.appendNamedFailures(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range iq.withNamedFiles {
		if err := iq.loadFiles(ctx, query, nodes,
			func(n *Incident) { n.appendNamedFiles(name) },
			func(n *Incident, e *File) { n.appendNamedFiles(name, e) }); err != nil {
			return nil, err
		}
	}
	for i := range iq.loadTotal {
		if err := iq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (iq *IncidentQuery) loadCounters(ctx context.Context, query *CounterQuery, nodes []*Incident, init func(*Incident), assign func(*Incident, *Counter)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Incident)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Counter(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(incident.CountersColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.incident_counters
		if fk == nil {
			return fmt.Errorf(`foreign-key "incident_counters" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "incident_counters" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (iq *IncidentQuery) loadStati(ctx context.Context, query *StatusQuery, nodes []*Incident, init func(*Incident), assign func(*Incident, *Status)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Incident)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Status(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(incident.StatiColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.incident_stati
		if fk == nil {
			return fmt.Errorf(`foreign-key "incident_stati" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "incident_stati" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (iq *IncidentQuery) loadFailures(ctx context.Context, query *FailureQuery, nodes []*Incident, init func(*Incident), assign func(*Incident, *Failure)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Incident)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Failure(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(incident.FailuresColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.incident_failures
		if fk == nil {
			return fmt.Errorf(`foreign-key "incident_failures" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "incident_failures" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (iq *IncidentQuery) loadFiles(ctx context.Context, query *FileQuery, nodes []*Incident, init func(*Incident), assign func(*Incident, *File)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Incident)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.File(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(incident.FilesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.incident_files
		if fk == nil {
			return fmt.Errorf(`foreign-key "incident_files" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "incident_files" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (iq *IncidentQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := iq.querySpec()
	if len(iq.modifiers) > 0 {
		_spec.Modifiers = iq.modifiers
	}
	_spec.Node.Columns = iq.ctx.Fields
	if len(iq.ctx.Fields) > 0 {
		_spec.Unique = iq.ctx.Unique != nil && *iq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, iq.driver, _spec)
}

func (iq *IncidentQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(incident.Table, incident.Columns, sqlgraph.NewFieldSpec(incident.FieldID, field.TypeInt))
	_spec.From = iq.sql
	if unique := iq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if iq.path != nil {
		_spec.Unique = true
	}
	if fields := iq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, incident.FieldID)
		for i := range fields {
			if fields[i] != incident.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := iq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := iq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := iq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := iq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (iq *IncidentQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(iq.driver.Dialect())
	t1 := builder.Table(incident.Table)
	columns := iq.ctx.Fields
	if len(columns) == 0 {
		columns = incident.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if iq.sql != nil {
		selector = iq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if iq.ctx.Unique != nil && *iq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range iq.predicates {
		p(selector)
	}
	for _, p := range iq.order {
		p(selector)
	}
	if offset := iq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := iq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WithNamedCounters tells the query-builder to eager-load the nodes that are connected to the "Counters"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithNamedCounters(name string, opts ...func(*CounterQuery)) *IncidentQuery {
	query := (&CounterClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if iq.withNamedCounters == nil {
		iq.withNamedCounters = make(map[string]*CounterQuery)
	}
	iq.withNamedCounters[name] = query
	return iq
}

// WithNamedStati tells the query-builder to eager-load the nodes that are connected to the "Stati"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithNamedStati(name string, opts ...func(*StatusQuery)) *IncidentQuery {
	query := (&StatusClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if iq.withNamedStati == nil {
		iq.withNamedStati = make(map[string]*StatusQuery)
	}
	iq.withNamedStati[name] = query
	return iq
}

// WithNamedFailures tells the query-builder to eager-load the nodes that are connected to the "Failures"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithNamedFailures(name string, opts ...func(*FailureQuery)) *IncidentQuery {
	query := (&FailureClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if iq.withNamedFailures == nil {
		iq.withNamedFailures = make(map[string]*FailureQuery)
	}
	iq.withNamedFailures[name] = query
	return iq
}

// WithNamedFiles tells the query-builder to eager-load the nodes that are connected to the "Files"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (iq *IncidentQuery) WithNamedFiles(name string, opts ...func(*FileQuery)) *IncidentQuery {
	query := (&FileClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if iq.withNamedFiles == nil {
		iq.withNamedFiles = make(map[string]*FileQuery)
	}
	iq.withNamedFiles[name] = query
	return iq
}

// IncidentGroupBy is the group-by builder for Incident entities.
type IncidentGroupBy struct {
	selector
	build *IncidentQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (igb *IncidentGroupBy) Aggregate(fns ...AggregateFunc) *IncidentGroupBy {
	igb.fns = append(igb.fns, fns...)
	return igb
}

// Scan applies the selector query and scans the result into the given value.
func (igb *IncidentGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, igb.build.ctx, ent.OpQueryGroupBy)
	if err := igb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*IncidentQuery, *IncidentGroupBy](ctx, igb.build, igb, igb.build.inters, v)
}

func (igb *IncidentGroupBy) sqlScan(ctx context.Context, root *IncidentQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(igb.fns))
	for _, fn := range igb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*igb.flds)+len(igb.fns))
		for _, f := range *igb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*igb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := igb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// IncidentSelect is the builder for selecting fields of Incident entities.
type IncidentSelect struct {
	*IncidentQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (is *IncidentSelect) Aggregate(fns ...AggregateFunc) *IncidentSelect {
	is.fns = append(is.fns, fns...)
	return is
}

// Scan applies the selector query and scans the result into the given value.
func (is *IncidentSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, is.ctx, ent.OpQuerySelect)
	if err := is.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*IncidentQuery, *IncidentSelect](ctx, is.IncidentQuery, is, is.inters, v)
}

func (is *IncidentSelect) sqlScan(ctx context.Context, root *IncidentQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(is.fns))
	for _, fn := range is.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*is.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := is.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/counter"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/status"
)

// AlertQuery is the builder for querying Alert entities.
type AlertQuery struct {
	config
	ctx               *QueryContext
	order             []alert.OrderOption
	inters            []Interceptor
	predicates        []predicate.Alert
	withCounters      *CounterQuery
	withStati         *StatusQuery
	withFailures      *FailureQuery
	withFiles         *FileQuery
	modifiers         []func(*sql.Selector)
	loadTotal         []func(context.Context, []*Alert) error
	withNamedCounters map[string]*CounterQuery
	withNamedStati    map[string]*StatusQuery
	withNamedFailures map[string]*FailureQuery
	withNamedFiles    map[string]*FileQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AlertQuery builder.
func (aq *AlertQuery) Where(ps ...predicate.Alert) *AlertQuery {
	aq.predicates = append(aq.predicates, ps...)
	return aq
}

// Limit the number of records to be returned by this query.
func (aq *AlertQuery) Limit(limit int) *AlertQuery {
	aq.ctx.Limit = &limit
	return aq
}

// Offset to start from.
func (aq *AlertQuery) Offset(offset int) *AlertQuery {
	aq.ctx.Offset = &offset
	return aq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (aq *AlertQuery) Unique(unique bool) *AlertQuery {
	aq.ctx.Unique = &unique
	return aq
}

// Order specifies how the records should be ordered.
func (aq *AlertQuery) Order(o ...alert.OrderOption) *AlertQuery {
	aq.order = append(aq.order, o...)
	return aq
}

// QueryCounters chains the current query on the "Counters" edge.
func (aq *AlertQuery) QueryCounters() *CounterQuery {
	query := (&CounterClient{config: aq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(alert.Table, alert.FieldID, selector),
			sqlgraph.To(counter.Table, counter.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, alert.CountersTable, alert.CountersColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryStati chains the current query on the "Stati" edge.
func (aq *AlertQuery) QueryStati() *StatusQuery {
	query := (&StatusClient{config: aq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(alert.Table, alert.FieldID, selector),
			sqlgraph.To(status.Table, status.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, alert.StatiTable, alert.StatiColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFailures chains the current query on the "Failures" edge.
func (aq *AlertQuery) QueryFailures() *FailureQuery {
	query := (&FailureClient{config: aq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(alert.Table, alert.FieldID, selector),
			sqlgraph.To(failure.Table, failure.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, alert.FailuresTable, alert.FailuresColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFiles chains the current query on the "Files" edge.
func (aq *AlertQuery) QueryFiles() *FileQuery {
	query := (&FileClient{config: aq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(alert.Table, alert.FieldID, selector),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, alert.FilesTable, alert.FilesColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Alert entity from the query.
// Returns a *NotFoundError when no Alert was found.
func (aq *AlertQuery) First(ctx context.Context) (*Alert, error) {
	nodes, err := aq.Limit(1).All(setContextOp(ctx, aq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{alert.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (aq *AlertQuery) FirstX(ctx context.Context) *Alert {
	node, err := aq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Alert ID from the query.
// Returns a *NotFoundError when no Alert ID was found.
func (aq *AlertQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = aq.Limit(1).IDs(setContextOp(ctx, aq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{alert.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (aq *AlertQuery) FirstIDX(ctx context.Context) int {
	id, err := aq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Alert entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Alert entity is found.
// Returns a *NotFoundError when no Alert entities are found.
func (aq *AlertQuery) Only(ctx context.Context) (*Alert, error) {
	nodes, err := aq.Limit(2).All(setContextOp(ctx, aq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{alert.Label}
	default:
		return nil, &NotSingularError{alert.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (aq *AlertQuery) OnlyX(ctx context.Context) *Alert {
	node, err := aq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Alert ID in the query.
// Returns a *NotSingularError when more than one Alert ID is found.
// Returns a *NotFoundError when no entities are found.
func (aq *AlertQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = aq.Limit(2).IDs(setContextOp(ctx, aq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{alert.Label}
	default:
		err = &NotSingularError{alert.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (aq *AlertQuery) OnlyIDX(ctx context.Context) int {
	id, err := aq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Alerts.
func (aq *AlertQuery) All(ctx context.Context) ([]*Alert, error) {
	ctx = setContextOp(ctx, aq.ctx, "All")
	if err := aq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Alert, *AlertQuery]()
	return withInterceptors[[]*Alert](ctx, aq, qr, aq.inters)
}

// AllX is like All, but panics if an error occurs.
func (aq *AlertQuery) AllX(ctx context.Context) []*Alert {
	nodes, err := aq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Alert IDs.
func (aq *AlertQuery) IDs(ctx context.Context) (ids []int, err error) {
	if aq.ctx.Unique == nil && aq.path != nil {
		aq.Unique(true)
	}
	ctx = setContextOp(ctx, aq.ctx, "IDs")
	if err = aq.Select(alert.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (aq *AlertQuery) IDsX(ctx context.Context) []int {
	ids, err := aq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (aq *AlertQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, aq.ctx, "Count")
	if err := aq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, aq, querierCount[*AlertQuery](), aq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (aq *AlertQuery) CountX(ctx context.Context) int {
	count, err := aq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (aq *AlertQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, aq.ctx, "Exist")
	switch _, err := aq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (aq *AlertQuery) ExistX(ctx context.Context) bool {
	exist, err := aq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AlertQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (aq *AlertQuery) Clone() *AlertQuery {
	if aq == nil {
		return nil
	}
	return &AlertQuery{
		config:       aq.config,
		ctx:          aq.ctx.Clone(),
		order:        append([]alert.OrderOption{}, aq.order...),
		inters:       append([]Interceptor{}, aq.inters...),
		predicates:   append([]predicate.Alert{}, aq.predicates...),
		withCounters: aq.withCounters.Clone(),
		withStati:    aq.withStati.Clone(),
		withFailures: aq.withFailures.Clone(),
		withFiles:    aq.withFiles.Clone(),
		// clone intermediate query.
		sql:  aq.sql.Clone(),
		path: aq.path,
	}
}

// WithCounters tells the query-builder to eager-load the nodes that are connected to
// the "Counters" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithCounters(opts ...func(*CounterQuery)) *AlertQuery {
	query := (&CounterClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aq.withCounters = query
	return aq
}

// WithStati tells the query-builder to eager-load the nodes that are connected to
// the "Stati" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithStati(opts ...func(*StatusQuery)) *AlertQuery {
	query := (&StatusClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aq.withStati = query
	return aq
}

// WithFailures tells the query-builder to eager-load the nodes that are connected to
// the "Failures" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithFailures(opts ...func(*FailureQuery)) *AlertQuery {
	query := (&FailureClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aq.withFailures = query
	return aq
}

// WithFiles tells the query-builder to eager-load the nodes that are connected to
// the "Files" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithFiles(opts ...func(*FileQuery)) *AlertQuery {
	query := (&FileClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aq.withFiles = query
	return aq
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
//	client.Alert.Query().
//		GroupBy(alert.FieldUUID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (aq *AlertQuery) GroupBy(field string, fields ...string) *AlertGroupBy {
	aq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AlertGroupBy{build: aq}
	grbuild.flds = &aq.ctx.Fields
	grbuild.label = alert.Label
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
//	client.Alert.Query().
//		Select(alert.FieldUUID).
//		Scan(ctx, &v)
func (aq *AlertQuery) Select(fields ...string) *AlertSelect {
	aq.ctx.Fields = append(aq.ctx.Fields, fields...)
	sbuild := &AlertSelect{AlertQuery: aq}
	sbuild.label = alert.Label
	sbuild.flds, sbuild.scan = &aq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AlertSelect configured with the given aggregations.
func (aq *AlertQuery) Aggregate(fns ...AggregateFunc) *AlertSelect {
	return aq.Select().Aggregate(fns...)
}

func (aq *AlertQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range aq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, aq); err != nil {
				return err
			}
		}
	}
	for _, f := range aq.ctx.Fields {
		if !alert.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if aq.path != nil {
		prev, err := aq.path(ctx)
		if err != nil {
			return err
		}
		aq.sql = prev
	}
	return nil
}

func (aq *AlertQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Alert, error) {
	var (
		nodes       = []*Alert{}
		_spec       = aq.querySpec()
		loadedTypes = [4]bool{
			aq.withCounters != nil,
			aq.withStati != nil,
			aq.withFailures != nil,
			aq.withFiles != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Alert).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Alert{config: aq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(aq.modifiers) > 0 {
		_spec.Modifiers = aq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, aq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := aq.withCounters; query != nil {
		if err := aq.loadCounters(ctx, query, nodes,
			func(n *Alert) { n.Edges.Counters = []*Counter{} },
			func(n *Alert, e *Counter) { n.Edges.Counters = append(n.Edges.Counters, e) }); err != nil {
			return nil, err
		}
	}
	if query := aq.withStati; query != nil {
		if err := aq.loadStati(ctx, query, nodes,
			func(n *Alert) { n.Edges.Stati = []*Status{} },
			func(n *Alert, e *Status) { n.Edges.Stati = append(n.Edges.Stati, e) }); err != nil {
			return nil, err
		}
	}
	if query := aq.withFailures; query != nil {
		if err := aq.loadFailures(ctx, query, nodes,
			func(n *Alert) { n.Edges.Failures = []*Failure{} },
			func(n *Alert, e *Failure) { n.Edges.Failures = append(n.Edges.Failures, e) }); err != nil {
			return nil, err
		}
	}
	if query := aq.withFiles; query != nil {
		if err := aq.loadFiles(ctx, query, nodes,
			func(n *Alert) { n.Edges.Files = []*File{} },
			func(n *Alert, e *File) { n.Edges.Files = append(n.Edges.Files, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range aq.withNamedCounters {
		if err := aq.loadCounters(ctx, query, nodes,
			func(n *Alert) { n.appendNamedCounters(name) },
			func(n *Alert, e *Counter) { n.appendNamedCounters(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range aq.withNamedStati {
		if err := aq.loadStati(ctx, query, nodes,
			func(n *Alert) { n.appendNamedStati(name) },
			func(n *Alert, e *Status) { n.appendNamedStati(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range aq.withNamedFailures {
		if err := aq.loadFailures(ctx, query, nodes,
			func(n *Alert) { n.appendNamedFailures(name) },
			func(n *Alert, e *Failure) { n.appendNamedFailures(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range aq.withNamedFiles {
		if err := aq.loadFiles(ctx, query, nodes,
			func(n *Alert) { n.appendNamedFiles(name) },
			func(n *Alert, e *File) { n.appendNamedFiles(name, e) }); err != nil {
			return nil, err
		}
	}
	for i := range aq.loadTotal {
		if err := aq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (aq *AlertQuery) loadCounters(ctx context.Context, query *CounterQuery, nodes []*Alert, init func(*Alert), assign func(*Alert, *Counter)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Alert)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Counter(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(alert.CountersColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.alert_counters
		if fk == nil {
			return fmt.Errorf(`foreign-key "alert_counters" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "alert_counters" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (aq *AlertQuery) loadStati(ctx context.Context, query *StatusQuery, nodes []*Alert, init func(*Alert), assign func(*Alert, *Status)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Alert)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Status(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(alert.StatiColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.alert_stati
		if fk == nil {
			return fmt.Errorf(`foreign-key "alert_stati" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "alert_stati" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (aq *AlertQuery) loadFailures(ctx context.Context, query *FailureQuery, nodes []*Alert, init func(*Alert), assign func(*Alert, *Failure)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Alert)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Failure(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(alert.FailuresColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.alert_failures
		if fk == nil {
			return fmt.Errorf(`foreign-key "alert_failures" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "alert_failures" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (aq *AlertQuery) loadFiles(ctx context.Context, query *FileQuery, nodes []*Alert, init func(*Alert), assign func(*Alert, *File)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Alert)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.File(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(alert.FilesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.alert_files
		if fk == nil {
			return fmt.Errorf(`foreign-key "alert_files" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "alert_files" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (aq *AlertQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := aq.querySpec()
	if len(aq.modifiers) > 0 {
		_spec.Modifiers = aq.modifiers
	}
	_spec.Node.Columns = aq.ctx.Fields
	if len(aq.ctx.Fields) > 0 {
		_spec.Unique = aq.ctx.Unique != nil && *aq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, aq.driver, _spec)
}

func (aq *AlertQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(alert.Table, alert.Columns, sqlgraph.NewFieldSpec(alert.FieldID, field.TypeInt))
	_spec.From = aq.sql
	if unique := aq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if aq.path != nil {
		_spec.Unique = true
	}
	if fields := aq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, alert.FieldID)
		for i := range fields {
			if fields[i] != alert.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := aq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := aq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := aq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := aq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (aq *AlertQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(aq.driver.Dialect())
	t1 := builder.Table(alert.Table)
	columns := aq.ctx.Fields
	if len(columns) == 0 {
		columns = alert.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if aq.sql != nil {
		selector = aq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if aq.ctx.Unique != nil && *aq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range aq.predicates {
		p(selector)
	}
	for _, p := range aq.order {
		p(selector)
	}
	if offset := aq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := aq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WithNamedCounters tells the query-builder to eager-load the nodes that are connected to the "Counters"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithNamedCounters(name string, opts ...func(*CounterQuery)) *AlertQuery {
	query := (&CounterClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if aq.withNamedCounters == nil {
		aq.withNamedCounters = make(map[string]*CounterQuery)
	}
	aq.withNamedCounters[name] = query
	return aq
}

// WithNamedStati tells the query-builder to eager-load the nodes that are connected to the "Stati"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithNamedStati(name string, opts ...func(*StatusQuery)) *AlertQuery {
	query := (&StatusClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if aq.withNamedStati == nil {
		aq.withNamedStati = make(map[string]*StatusQuery)
	}
	aq.withNamedStati[name] = query
	return aq
}

// WithNamedFailures tells the query-builder to eager-load the nodes that are connected to the "Failures"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithNamedFailures(name string, opts ...func(*FailureQuery)) *AlertQuery {
	query := (&FailureClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if aq.withNamedFailures == nil {
		aq.withNamedFailures = make(map[string]*FailureQuery)
	}
	aq.withNamedFailures[name] = query
	return aq
}

// WithNamedFiles tells the query-builder to eager-load the nodes that are connected to the "Files"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (aq *AlertQuery) WithNamedFiles(name string, opts ...func(*FileQuery)) *AlertQuery {
	query := (&FileClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if aq.withNamedFiles == nil {
		aq.withNamedFiles = make(map[string]*FileQuery)
	}
	aq.withNamedFiles[name] = query
	return aq
}

// AlertGroupBy is the group-by builder for Alert entities.
type AlertGroupBy struct {
	selector
	build *AlertQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (agb *AlertGroupBy) Aggregate(fns ...AggregateFunc) *AlertGroupBy {
	agb.fns = append(agb.fns, fns...)
	return agb
}

// Scan applies the selector query and scans the result into the given value.
func (agb *AlertGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, agb.build.ctx, "GroupBy")
	if err := agb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AlertQuery, *AlertGroupBy](ctx, agb.build, agb, agb.build.inters, v)
}

func (agb *AlertGroupBy) sqlScan(ctx context.Context, root *AlertQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(agb.fns))
	for _, fn := range agb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*agb.flds)+len(agb.fns))
		for _, f := range *agb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*agb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := agb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AlertSelect is the builder for selecting fields of Alert entities.
type AlertSelect struct {
	*AlertQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (as *AlertSelect) Aggregate(fns ...AggregateFunc) *AlertSelect {
	as.fns = append(as.fns, fns...)
	return as
}

// Scan applies the selector query and scans the result into the given value.
func (as *AlertSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, as.ctx, "Select")
	if err := as.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AlertQuery, *AlertSelect](ctx, as.AlertQuery, as, as.inters, v)
}

func (as *AlertSelect) sqlScan(ctx context.Context, root *AlertQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(as.fns))
	for _, fn := range as.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*as.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := as.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/99designs/gqlgen/graphql"
	"github.com/hashicorp/go-multierror"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/alert"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/counter"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/file"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/incident"
	"github.com/vogtp/som/pkg/visualiser/webstatus/ent/ent/status"
	"golang.org/x/sync/semaphore"
)

// Noder wraps the basic Node method.
type Noder interface {
	Node(context.Context) (*Node, error)
}

// Node in the graph.
type Node struct {
	ID     int      `json:"id,omitempty"`     // node id.
	Type   string   `json:"type,omitempty"`   // node type.
	Fields []*Field `json:"fields,omitempty"` // node fields.
	Edges  []*Edge  `json:"edges,omitempty"`  // node edges.
}

// Field of a node.
type Field struct {
	Type  string `json:"type,omitempty"`  // field type.
	Name  string `json:"name,omitempty"`  // field name (as in struct).
	Value string `json:"value,omitempty"` // stringified value.
}

// Edges between two nodes.
type Edge struct {
	Type string `json:"type,omitempty"` // edge type.
	Name string `json:"name,omitempty"` // edge name.
	IDs  []int  `json:"ids,omitempty"`  // node ids (where this edge point to).
}

func (a *Alert) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     a.ID,
		Type:   "Alert",
		Fields: make([]*Field, 10),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(a.Level); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "int",
		Name:  "Level",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.UUID); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "uuid.UUID",
		Name:  "UUID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.IncidentID); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "uuid.UUID",
		Name:  "IncidentID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.Time); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "time.Time",
		Name:  "Time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.Username); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "Username",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.Region); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "Region",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.ProbeOS); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "ProbeOS",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.ProbeHost); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "string",
		Name:  "ProbeHost",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.Error); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "Error",
		Value: string(buf),
	}
	return node, nil
}

func (c *Counter) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     c.ID,
		Type:   "Counter",
		Fields: make([]*Field, 2),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(c.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Value); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "Value",
		Value: string(buf),
	}
	return node, nil
}

func (f *Failure) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     f.ID,
		Type:   "Failure",
		Fields: make([]*Field, 2),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(f.Error); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "Error",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Idx); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "int",
		Name:  "Idx",
		Value: string(buf),
	}
	return node, nil
}

func (f *File) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     f.ID,
		Type:   "File",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(f.UUID); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "uuid.UUID",
		Name:  "UUID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Name); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Type); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Ext); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Ext",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Size); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "Size",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Payload); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "[]byte",
		Name:  "payload",
		Value: string(buf),
	}
	return node, nil
}

func (i *Incident) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     i.ID,
		Type:   "Incident",
		Fields: make([]*Field, 13),
		Edges:  make([]*Edge, 4),
	}
	var buf []byte
	if buf, err = json.Marshal(i.Level); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "int",
		Name:  "Level",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.Start); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "Start",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.End); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "time.Time",
		Name:  "End",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.State); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "[]byte",
		Name:  "State",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.UUID); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "uuid.UUID",
		Name:  "UUID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.IncidentID); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "uuid.UUID",
		Name:  "IncidentID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.Name); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.Time); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "time.Time",
		Name:  "Time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.Username); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "string",
		Name:  "Username",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.Region); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "Region",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.ProbeOS); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "string",
		Name:  "ProbeOS",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.ProbeHost); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "string",
		Name:  "ProbeHost",
		Value: string(buf),
	}
	if buf, err = json.Marshal(i.Error); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "string",
		Name:  "Error",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "Counter",
		Name: "Counters",
	}
	err = i.QueryCounters().
		Select(counter.FieldID).
		Scan(ctx, &node.Edges[0].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "Status",
		Name: "Stati",
	}
	err = i.QueryStati().
		Select(status.FieldID).
		Scan(ctx, &node.Edges[1].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		Type: "Failure",
		Name: "Failures",
	}
	err = i.QueryFailures().
		Select(failure.FieldID).
		Scan(ctx, &node.Edges[2].IDs)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		Type: "File",
		Name: "Files",
	}
	err = i.QueryFiles().
		Select(file.FieldID).
		Scan(ctx, &node.Edges[3].IDs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (s *Status) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     s.ID,
		Type:   "Status",
		Fields: make([]*Field, 2),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(s.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.Value); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "Value",
		Value: string(buf),
	}
	return node, nil
}

func (c *Client) Node(ctx context.Context, id int) (*Node, error) {
	n, err := c.Noder(ctx, id)
	if err != nil {
		return nil, err
	}
	return n.Node(ctx)
}

var errNodeInvalidID = &NotFoundError{"node"}

// NodeOption allows configuring the Noder execution using functional options.
type NodeOption func(*nodeOptions)

// WithNodeType sets the node Type resolver function (i.e. the table to query).
// If was not provided, the table will be derived from the universal-id
// configuration as described in: https://entgo.io/docs/migrate/#universal-ids.
func WithNodeType(f func(context.Context, int) (string, error)) NodeOption {
	return func(o *nodeOptions) {
		o.nodeType = f
	}
}

// WithFixedNodeType sets the Type of the node to a fixed value.
func WithFixedNodeType(t string) NodeOption {
	return WithNodeType(func(context.Context, int) (string, error) {
		return t, nil
	})
}

type nodeOptions struct {
	nodeType func(context.Context, int) (string, error)
}

func (c *Client) newNodeOpts(opts []NodeOption) *nodeOptions {
	nopts := &nodeOptions{}
	for _, opt := range opts {
		opt(nopts)
	}
	if nopts.nodeType == nil {
		nopts.nodeType = func(ctx context.Context, id int) (string, error) {
			return c.tables.nodeType(ctx, c.driver, id)
		}
	}
	return nopts
}

// Noder returns a Node by its id. If the NodeType was not provided, it will
// be derived from the id value according to the universal-id configuration.
//
//	c.Noder(ctx, id)
//	c.Noder(ctx, id, ent.WithNodeType(typeResolver))
func (c *Client) Noder(ctx context.Context, id int, opts ...NodeOption) (_ Noder, err error) {
	defer func() {
		if IsNotFound(err) {
			err = multierror.Append(err, entgql.ErrNodeNotFound(id))
		}
	}()
	table, err := c.newNodeOpts(opts).nodeType(ctx, id)
	if err != nil {
		return nil, err
	}
	return c.noder(ctx, table, id)
}

func (c *Client) noder(ctx context.Context, table string, id int) (Noder, error) {
	switch table {
	case alert.Table:
		query := c.Alert.Query().
			Where(alert.ID(id))
		query, err := query.CollectFields(ctx, "Alert")
		if err != nil {
			return nil, err
		}
		n, err := query.Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case counter.Table:
		query := c.Counter.Query().
			Where(counter.ID(id))
		query, err := query.CollectFields(ctx, "Counter")
		if err != nil {
			return nil, err
		}
		n, err := query.Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case failure.Table:
		query := c.Failure.Query().
			Where(failure.ID(id))
		query, err := query.CollectFields(ctx, "Failure")
		if err != nil {
			return nil, err
		}
		n, err := query.Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case file.Table:
		query := c.File.Query().
			Where(file.ID(id))
		query, err := query.CollectFields(ctx, "File")
		if err != nil {
			return nil, err
		}
		n, err := query.Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case incident.Table:
		query := c.Incident.Query().
			Where(incident.ID(id))
		query, err := query.CollectFields(ctx, "Incident")
		if err != nil {
			return nil, err
		}
		n, err := query.Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case status.Table:
		query := c.Status.Query().
			Where(status.ID(id))
		query, err := query.CollectFields(ctx, "Status")
		if err != nil {
			return nil, err
		}
		n, err := query.Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	default:
		return nil, fmt.Errorf("cannot resolve noder from table %q: %w", table, errNodeInvalidID)
	}
}

func (c *Client) Noders(ctx context.Context, ids []int, opts ...NodeOption) ([]Noder, error) {
	switch len(ids) {
	case 1:
		noder, err := c.Noder(ctx, ids[0], opts...)
		if err != nil {
			return nil, err
		}
		return []Noder{noder}, nil
	case 0:
		return []Noder{}, nil
	}

	noders := make([]Noder, len(ids))
	errors := make([]error, len(ids))
	tables := make(map[string][]int)
	id2idx := make(map[int][]int, len(ids))
	nopts := c.newNodeOpts(opts)
	for i, id := range ids {
		table, err := nopts.nodeType(ctx, id)
		if err != nil {
			errors[i] = err
			continue
		}
		tables[table] = append(tables[table], id)
		id2idx[id] = append(id2idx[id], i)
	}

	for table, ids := range tables {
		nodes, err := c.noders(ctx, table, ids)
		if err != nil {
			for _, id := range ids {
				for _, idx := range id2idx[id] {
					errors[idx] = err
				}
			}
		} else {
			for i, id := range ids {
				for _, idx := range id2idx[id] {
					noders[idx] = nodes[i]
				}
			}
		}
	}

	for i, id := range ids {
		if errors[i] == nil {
			if noders[i] != nil {
				continue
			}
			errors[i] = entgql.ErrNodeNotFound(id)
		} else if IsNotFound(errors[i]) {
			errors[i] = multierror.Append(errors[i], entgql.ErrNodeNotFound(id))
		}
		ctx := graphql.WithPathContext(ctx,
			graphql.NewPathWithIndex(i),
		)
		graphql.AddError(ctx, errors[i])
	}
	return noders, nil
}

func (c *Client) noders(ctx context.Context, table string, ids []int) ([]Noder, error) {
	noders := make([]Noder, len(ids))
	idmap := make(map[int][]*Noder, len(ids))
	for i, id := range ids {
		idmap[id] = append(idmap[id], &noders[i])
	}
	switch table {
	case alert.Table:
		query := c.Alert.Query().
			Where(alert.IDIn(ids...))
		query, err := query.CollectFields(ctx, "Alert")
		if err != nil {
			return nil, err
		}
		nodes, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			for _, noder := range idmap[node.ID] {
				*noder = node
			}
		}
	case counter.Table:
		query := c.Counter.Query().
			Where(counter.IDIn(ids...))
		query, err := query.CollectFields(ctx, "Counter")
		if err != nil {
			return nil, err
		}
		nodes, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			for _, noder := range idmap[node.ID] {
				*noder = node
			}
		}
	case failure.Table:
		query := c.Failure.Query().
			Where(failure.IDIn(ids...))
		query, err := query.CollectFields(ctx, "Failure")
		if err != nil {
			return nil, err
		}
		nodes, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			for _, noder := range idmap[node.ID] {
				*noder = node
			}
		}
	case file.Table:
		query := c.File.Query().
			Where(file.IDIn(ids...))
		query, err := query.CollectFields(ctx, "File")
		if err != nil {
			return nil, err
		}
		nodes, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			for _, noder := range idmap[node.ID] {
				*noder = node
			}
		}
	case incident.Table:
		query := c.Incident.Query().
			Where(incident.IDIn(ids...))
		query, err := query.CollectFields(ctx, "Incident")
		if err != nil {
			return nil, err
		}
		nodes, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			for _, noder := range idmap[node.ID] {
				*noder = node
			}
		}
	case status.Table:
		query := c.Status.Query().
			Where(status.IDIn(ids...))
		query, err := query.CollectFields(ctx, "Status")
		if err != nil {
			return nil, err
		}
		nodes, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			for _, noder := range idmap[node.ID] {
				*noder = node
			}
		}
	default:
		return nil, fmt.Errorf("cannot resolve noders from table %q: %w", table, errNodeInvalidID)
	}
	return noders, nil
}

type tables struct {
	once  sync.Once
	sem   *semaphore.Weighted
	value atomic.Value
}

func (t *tables) nodeType(ctx context.Context, drv dialect.Driver, id int) (string, error) {
	tables, err := t.Load(ctx, drv)
	if err != nil {
		return "", err
	}
	idx := int(id / (1<<32 - 1))
	if idx < 0 || idx >= len(tables) {
		return "", fmt.Errorf("cannot resolve table from id %v: %w", id, errNodeInvalidID)
	}
	return tables[idx], nil
}

func (t *tables) Load(ctx context.Context, drv dialect.Driver) ([]string, error) {
	if tables := t.value.Load(); tables != nil {
		return tables.([]string), nil
	}
	t.once.Do(func() { t.sem = semaphore.NewWeighted(1) })
	if err := t.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer t.sem.Release(1)
	if tables := t.value.Load(); tables != nil {
		return tables.([]string), nil
	}
	tables, err := t.load(ctx, drv)
	if err == nil {
		t.value.Store(tables)
	}
	return tables, err
}

func (*tables) load(ctx context.Context, drv dialect.Driver) ([]string, error) {
	rows := &sql.Rows{}
	query, args := sql.Dialect(drv.Dialect()).
		Select("type").
		From(sql.Table(schema.TypeTable)).
		OrderBy(sql.Asc("id")).
		Query()
	if err := drv.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	return tables, sql.ScanSlice(rows, &tables)
}

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

// Incident is the model entity for the Incident schema.
type Incident struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// UUID holds the value of the "UUID" field.
	UUID uuid.UUID `json:"UUID,omitempty"`
	// IncidentID holds the value of the "IncidentID" field.
	IncidentID uuid.UUID `json:"IncidentID,omitempty"`
	// Name holds the value of the "Name" field.
	Name string `json:"Name,omitempty"`
	// Time holds the value of the "Time" field.
	Time time.Time `json:"Time,omitempty"`
	// IntLevel holds the value of the "IntLevel" field.
	IntLevel int `json:"IntLevel,omitempty"`
	// Username holds the value of the "Username" field.
	Username string `json:"Username,omitempty"`
	// Region holds the value of the "Region" field.
	Region string `json:"Region,omitempty"`
	// ProbeOS holds the value of the "ProbeOS" field.
	ProbeOS string `json:"ProbeOS,omitempty"`
	// ProbeHost holds the value of the "ProbeHost" field.
	ProbeHost string `json:"ProbeHost,omitempty"`
	// Error holds the value of the "Error" field.
	Error string `json:"Error,omitempty"`
	// Start holds the value of the "Start" field.
	Start time.Time `json:"Start,omitempty"`
	// End holds the value of the "End" field.
	End time.Time `json:"End,omitempty"`
	// State holds the value of the "State" field.
	State []byte `json:"State,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the IncidentQuery when eager-loading is set.
	Edges        IncidentEdges `json:"edges"`
	selectValues sql.SelectValues
}

// IncidentEdges holds the relations/edges for other nodes in the graph.
type IncidentEdges struct {
	// Counters holds the value of the Counters edge.
	Counters []*Counter `json:"Counters,omitempty"`
	// Stati holds the value of the Stati edge.
	Stati []*Status `json:"Stati,omitempty"`
	// Failures holds the value of the Failures edge.
	Failures []*Failure `json:"Failures,omitempty"`
	// Files holds the value of the Files edge.
	Files []*File `json:"Files,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
	// totalCount holds the count of the edges above.
	totalCount [4]map[string]int

	namedCounters map[string][]*Counter
	namedStati    map[string][]*Status
	namedFailures map[string][]*Failure
	namedFiles    map[string][]*File
}

// CountersOrErr returns the Counters value or an error if the edge
// was not loaded in eager-loading.
func (e IncidentEdges) CountersOrErr() ([]*Counter, error) {
	if e.loadedTypes[0] {
		return e.Counters, nil
	}
	return nil, &NotLoadedError{edge: "Counters"}
}

// StatiOrErr returns the Stati value or an error if the edge
// was not loaded in eager-loading.
func (e IncidentEdges) StatiOrErr() ([]*Status, error) {
	if e.loadedTypes[1] {
		return e.Stati, nil
	}
	return nil, &NotLoadedError{edge: "Stati"}
}

// FailuresOrErr returns the Failures value or an error if the edge
// was not loaded in eager-loading.
func (e IncidentEdges) FailuresOrErr() ([]*Failure, error) {
	if e.loadedTypes[2] {
		return e.Failures, nil
	}
	return nil, &NotLoadedError{edge: "Failures"}
}

// FilesOrErr returns the Files value or an error if the edge
// was not loaded in eager-loading.
func (e IncidentEdges) FilesOrErr() ([]*File, error) {
	if e.loadedTypes[3] {
		return e.Files, nil
	}
	return nil, &NotLoadedError{edge: "Files"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Incident) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case incident.FieldState:
			values[i] = new([]byte)
		case incident.FieldID, incident.FieldIntLevel:
			values[i] = new(sql.NullInt64)
		case incident.FieldName, incident.FieldUsername, incident.FieldRegion, incident.FieldProbeOS, incident.FieldProbeHost, incident.FieldError:
			values[i] = new(sql.NullString)
		case incident.FieldTime, incident.FieldStart, incident.FieldEnd:
			values[i] = new(sql.NullTime)
		case incident.FieldUUID, incident.FieldIncidentID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Incident fields.
func (i *Incident) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case incident.FieldID:
			value, ok := values[j].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			i.ID = int(value.Int64)
		case incident.FieldUUID:
			if value, ok := values[j].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field UUID", values[j])
			} else if value != nil {
				i.UUID = *value
			}
		case incident.FieldIncidentID:
			if value, ok := values[j].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field IncidentID", values[j])
			} else if value != nil {
				i.IncidentID = *value
			}
		case incident.FieldName:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field Name", values[j])
			} else if value.Valid {
				i.Name = value.String
			}
		case incident.FieldTime:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field Time", values[j])
			} else if value.Valid {
				i.Time = value.Time
			}
		case incident.FieldIntLevel:
			if value, ok := values[j].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field IntLevel", values[j])
			} else if value.Valid {
				i.IntLevel = int(value.Int64)
			}
		case incident.FieldUsername:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field Username", values[j])
			} else if value.Valid {
				i.Username = value.String
			}
		case incident.FieldRegion:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field Region", values[j])
			} else if value.Valid {
				i.Region = value.String
			}
		case incident.FieldProbeOS:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ProbeOS", values[j])
			} else if value.Valid {
				i.ProbeOS = value.String
			}
		case incident.FieldProbeHost:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ProbeHost", values[j])
			} else if value.Valid {
				i.ProbeHost = value.String
			}
		case incident.FieldError:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field Error", values[j])
			} else if value.Valid {
				i.Error = value.String
			}
		case incident.FieldStart:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field Start", values[j])
			} else if value.Valid {
				i.Start = value.Time
			}
		case incident.FieldEnd:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field End", values[j])
			} else if value.Valid {
				i.End = value.Time
			}
		case incident.FieldState:
			if value, ok := values[j].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field State", values[j])
			} else if value != nil {
				i.State = *value
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Incident.
// This includes values selected through modifiers, order, etc.
func (i *Incident) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// QueryCounters queries the "Counters" edge of the Incident entity.
func (i *Incident) QueryCounters() *CounterQuery {
	return NewIncidentClient(i.config).QueryCounters(i)
}

// QueryStati queries the "Stati" edge of the Incident entity.
func (i *Incident) QueryStati() *StatusQuery {
	return NewIncidentClient(i.config).QueryStati(i)
}

// QueryFailures queries the "Failures" edge of the Incident entity.
func (i *Incident) QueryFailures() *FailureQuery {
	return NewIncidentClient(i.config).QueryFailures(i)
}

// QueryFiles queries the "Files" edge of the Incident entity.
func (i *Incident) QueryFiles() *FileQuery {
	return NewIncidentClient(i.config).QueryFiles(i)
}

// Update returns a builder for updating this Incident.
// Note that you need to call Incident.Unwrap() before calling this method if this Incident
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Incident) Update() *IncidentUpdateOne {
	return NewIncidentClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Incident entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Incident) Unwrap() *Incident {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Incident is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Incident) String() string {
	var builder strings.Builder
	builder.WriteString("Incident(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("UUID=")
	builder.WriteString(fmt.Sprintf("%v", i.UUID))
	builder.WriteString(", ")
	builder.WriteString("IncidentID=")
	builder.WriteString(fmt.Sprintf("%v", i.IncidentID))
	builder.WriteString(", ")
	builder.WriteString("Name=")
	builder.WriteString(i.Name)
	builder.WriteString(", ")
	builder.WriteString("Time=")
	builder.WriteString(i.Time.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("IntLevel=")
	builder.WriteString(fmt.Sprintf("%v", i.IntLevel))
	builder.WriteString(", ")
	builder.WriteString("Username=")
	builder.WriteString(i.Username)
	builder.WriteString(", ")
	builder.WriteString("Region=")
	builder.WriteString(i.Region)
	builder.WriteString(", ")
	builder.WriteString("ProbeOS=")
	builder.WriteString(i.ProbeOS)
	builder.WriteString(", ")
	builder.WriteString("ProbeHost=")
	builder.WriteString(i.ProbeHost)
	builder.WriteString(", ")
	builder.WriteString("Error=")
	builder.WriteString(i.Error)
	builder.WriteString(", ")
	builder.WriteString("Start=")
	builder.WriteString(i.Start.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("End=")
	builder.WriteString(i.End.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("State=")
	builder.WriteString(fmt.Sprintf("%v", i.State))
	builder.WriteByte(')')
	return builder.String()
}

// NamedCounters returns the Counters named value or an error if the edge was not
// loaded in eager-loading with this name.
func (i *Incident) NamedCounters(name string) ([]*Counter, error) {
	if i.Edges.namedCounters == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := i.Edges.namedCounters[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (i *Incident) appendNamedCounters(name string, edges ...*Counter) {
	if i.Edges.namedCounters == nil {
		i.Edges.namedCounters = make(map[string][]*Counter)
	}
	if len(edges) == 0 {
		i.Edges.namedCounters[name] = []*Counter{}
	} else {
		i.Edges.namedCounters[name] = append(i.Edges.namedCounters[name], edges...)
	}
}

// NamedStati returns the Stati named value or an error if the edge was not
// loaded in eager-loading with this name.
func (i *Incident) NamedStati(name string) ([]*Status, error) {
	if i.Edges.namedStati == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := i.Edges.namedStati[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (i *Incident) appendNamedStati(name string, edges ...*Status) {
	if i.Edges.namedStati == nil {
		i.Edges.namedStati = make(map[string][]*Status)
	}
	if len(edges) == 0 {
		i.Edges.namedStati[name] = []*Status{}
	} else {
		i.Edges.namedStati[name] = append(i.Edges.namedStati[name], edges...)
	}
}

// NamedFailures returns the Failures named value or an error if the edge was not
// loaded in eager-loading with this name.
func (i *Incident) NamedFailures(name string) ([]*Failure, error) {
	if i.Edges.namedFailures == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := i.Edges.namedFailures[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (i *Incident) appendNamedFailures(name string, edges ...*Failure) {
	if i.Edges.namedFailures == nil {
		i.Edges.namedFailures = make(map[string][]*Failure)
	}
	if len(edges) == 0 {
		i.Edges.namedFailures[name] = []*Failure{}
	} else {
		i.Edges.namedFailures[name] = append(i.Edges.namedFailures[name], edges...)
	}
}

// NamedFiles returns the Files named value or an error if the edge was not
// loaded in eager-loading with this name.
func (i *Incident) NamedFiles(name string) ([]*File, error) {
	if i.Edges.namedFiles == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := i.Edges.namedFiles[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (i *Incident) appendNamedFiles(name string, edges ...*File) {
	if i.Edges.namedFiles == nil {
		i.Edges.namedFiles = make(map[string][]*File)
	}
	if len(edges) == 0 {
		i.Edges.namedFiles[name] = []*File{}
	} else {
		i.Edges.namedFiles[name] = append(i.Edges.namedFiles[name], edges...)
	}
}

// Incidents is a parsable slice of Incident.
type Incidents []*Incident

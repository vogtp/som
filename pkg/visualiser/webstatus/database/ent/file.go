// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent/file"
)

// File is the model entity for the File schema.
type File struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// UUID holds the value of the "UUID" field.
	UUID uuid.UUID `json:"UUID,omitempty"`
	// Name holds the value of the "Name" field.
	Name string `json:"Name,omitempty"`
	// Type holds the value of the "Type" field.
	Type string `json:"Type,omitempty"`
	// Ext holds the value of the "Ext" field.
	Ext string `json:"Ext,omitempty"`
	// Size holds the value of the "Size" field.
	Size int `json:"Size,omitempty"`
	// Payload holds the value of the "payload" field.
	Payload        []byte `json:"payload,omitempty"`
	alert_files    *int
	incident_files *int
}

// scanValues returns the types for scanning values from sql.Rows.
func (*File) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case file.FieldPayload:
			values[i] = new([]byte)
		case file.FieldID, file.FieldSize:
			values[i] = new(sql.NullInt64)
		case file.FieldName, file.FieldType, file.FieldExt:
			values[i] = new(sql.NullString)
		case file.FieldUUID:
			values[i] = new(uuid.UUID)
		case file.ForeignKeys[0]: // alert_files
			values[i] = new(sql.NullInt64)
		case file.ForeignKeys[1]: // incident_files
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type File", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the File fields.
func (f *File) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case file.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			f.ID = int(value.Int64)
		case file.FieldUUID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field UUID", values[i])
			} else if value != nil {
				f.UUID = *value
			}
		case file.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field Name", values[i])
			} else if value.Valid {
				f.Name = value.String
			}
		case file.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field Type", values[i])
			} else if value.Valid {
				f.Type = value.String
			}
		case file.FieldExt:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field Ext", values[i])
			} else if value.Valid {
				f.Ext = value.String
			}
		case file.FieldSize:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field Size", values[i])
			} else if value.Valid {
				f.Size = int(value.Int64)
			}
		case file.FieldPayload:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field payload", values[i])
			} else if value != nil {
				f.Payload = *value
			}
		case file.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field alert_files", value)
			} else if value.Valid {
				f.alert_files = new(int)
				*f.alert_files = int(value.Int64)
			}
		case file.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field incident_files", value)
			} else if value.Valid {
				f.incident_files = new(int)
				*f.incident_files = int(value.Int64)
			}
		}
	}
	return nil
}

// Update returns a builder for updating this File.
// Note that you need to call File.Unwrap() before calling this method if this File
// was returned from a transaction, and the transaction was committed or rolled back.
func (f *File) Update() *FileUpdateOne {
	return (&FileClient{config: f.config}).UpdateOne(f)
}

// Unwrap unwraps the File entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (f *File) Unwrap() *File {
	_tx, ok := f.config.driver.(*txDriver)
	if !ok {
		panic("ent: File is not a transactional entity")
	}
	f.config.driver = _tx.drv
	return f
}

// String implements the fmt.Stringer.
func (f *File) String() string {
	var builder strings.Builder
	builder.WriteString("File(")
	builder.WriteString(fmt.Sprintf("id=%v, ", f.ID))
	builder.WriteString("UUID=")
	builder.WriteString(fmt.Sprintf("%v", f.UUID))
	builder.WriteString(", ")
	builder.WriteString("Name=")
	builder.WriteString(f.Name)
	builder.WriteString(", ")
	builder.WriteString("Type=")
	builder.WriteString(f.Type)
	builder.WriteString(", ")
	builder.WriteString("Ext=")
	builder.WriteString(f.Ext)
	builder.WriteString(", ")
	builder.WriteString("Size=")
	builder.WriteString(fmt.Sprintf("%v", f.Size))
	builder.WriteString(", ")
	builder.WriteString("payload=")
	builder.WriteString(fmt.Sprintf("%v", f.Payload))
	builder.WriteByte(')')
	return builder.String()
}

// Files is a parsable slice of File.
type Files []*File

func (f Files) config(cfg config) {
	for _i := range f {
		f[_i].config = cfg
	}
}

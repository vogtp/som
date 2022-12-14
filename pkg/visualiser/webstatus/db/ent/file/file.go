// Code generated by ent, DO NOT EDIT.

package file

const (
	// Label holds the string label denoting the file type in the database.
	Label = "file"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUUID holds the string denoting the uuid field in the database.
	FieldUUID = "uuid"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldExt holds the string denoting the ext field in the database.
	FieldExt = "ext"
	// FieldSize holds the string denoting the size field in the database.
	FieldSize = "size"
	// FieldPayload holds the string denoting the payload field in the database.
	FieldPayload = "payload"
	// Table holds the table name of the file in the database.
	Table = "files"
)

// Columns holds all SQL columns for file fields.
var Columns = []string{
	FieldID,
	FieldUUID,
	FieldName,
	FieldType,
	FieldExt,
	FieldSize,
	FieldPayload,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "files"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"alert_files",
	"incident_files",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

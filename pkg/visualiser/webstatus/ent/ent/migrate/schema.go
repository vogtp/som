// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AlertsColumns holds the columns for the "alerts" table.
	AlertsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "level", Type: field.TypeInt},
		{Name: "uuid", Type: field.TypeUUID, Unique: true},
		{Name: "incident_id", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString},
		{Name: "time", Type: field.TypeTime},
		{Name: "username", Type: field.TypeString},
		{Name: "region", Type: field.TypeString},
		{Name: "probe_os", Type: field.TypeString},
		{Name: "probe_host", Type: field.TypeString},
		{Name: "error", Type: field.TypeString, Nullable: true},
	}
	// AlertsTable holds the schema information for the "alerts" table.
	AlertsTable = &schema.Table{
		Name:       "alerts",
		Columns:    AlertsColumns,
		PrimaryKey: []*schema.Column{AlertsColumns[0]},
	}
	// CountersColumns holds the columns for the "counters" table.
	CountersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "value", Type: field.TypeString},
		{Name: "alert_counters", Type: field.TypeInt, Nullable: true},
		{Name: "incident_counters", Type: field.TypeInt, Nullable: true},
	}
	// CountersTable holds the schema information for the "counters" table.
	CountersTable = &schema.Table{
		Name:       "counters",
		Columns:    CountersColumns,
		PrimaryKey: []*schema.Column{CountersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "counters_alerts_Counters",
				Columns:    []*schema.Column{CountersColumns[3]},
				RefColumns: []*schema.Column{AlertsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "counters_incidents_Counters",
				Columns:    []*schema.Column{CountersColumns[4]},
				RefColumns: []*schema.Column{IncidentsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// FailuresColumns holds the columns for the "failures" table.
	FailuresColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "error", Type: field.TypeString},
		{Name: "idx", Type: field.TypeInt},
		{Name: "alert_failures", Type: field.TypeInt, Nullable: true},
		{Name: "incident_failures", Type: field.TypeInt, Nullable: true},
	}
	// FailuresTable holds the schema information for the "failures" table.
	FailuresTable = &schema.Table{
		Name:       "failures",
		Columns:    FailuresColumns,
		PrimaryKey: []*schema.Column{FailuresColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "failures_alerts_Failures",
				Columns:    []*schema.Column{FailuresColumns[3]},
				RefColumns: []*schema.Column{AlertsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "failures_incidents_Failures",
				Columns:    []*schema.Column{FailuresColumns[4]},
				RefColumns: []*schema.Column{IncidentsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// FilesColumns holds the columns for the "files" table.
	FilesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "uuid", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "ext", Type: field.TypeString},
		{Name: "size", Type: field.TypeInt},
		{Name: "payload", Type: field.TypeBytes},
		{Name: "alert_files", Type: field.TypeInt, Nullable: true},
		{Name: "incident_files", Type: field.TypeInt, Nullable: true},
	}
	// FilesTable holds the schema information for the "files" table.
	FilesTable = &schema.Table{
		Name:       "files",
		Columns:    FilesColumns,
		PrimaryKey: []*schema.Column{FilesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "files_alerts_Files",
				Columns:    []*schema.Column{FilesColumns[7]},
				RefColumns: []*schema.Column{AlertsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "files_incidents_Files",
				Columns:    []*schema.Column{FilesColumns[8]},
				RefColumns: []*schema.Column{IncidentsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// IncidentsColumns holds the columns for the "incidents" table.
	IncidentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "level", Type: field.TypeInt},
		{Name: "start", Type: field.TypeTime},
		{Name: "end", Type: field.TypeTime},
		{Name: "state", Type: field.TypeBytes},
		{Name: "uuid", Type: field.TypeUUID, Unique: true},
		{Name: "incident_id", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString},
		{Name: "time", Type: field.TypeTime},
		{Name: "username", Type: field.TypeString},
		{Name: "region", Type: field.TypeString},
		{Name: "probe_os", Type: field.TypeString},
		{Name: "probe_host", Type: field.TypeString},
		{Name: "error", Type: field.TypeString, Nullable: true},
	}
	// IncidentsTable holds the schema information for the "incidents" table.
	IncidentsTable = &schema.Table{
		Name:       "incidents",
		Columns:    IncidentsColumns,
		PrimaryKey: []*schema.Column{IncidentsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "incident_start",
				Unique:  false,
				Columns: []*schema.Column{IncidentsColumns[2]},
			},
			{
				Name:    "incident_end",
				Unique:  false,
				Columns: []*schema.Column{IncidentsColumns[3]},
			},
			{
				Name:    "incident_uuid",
				Unique:  true,
				Columns: []*schema.Column{IncidentsColumns[5]},
			},
			{
				Name:    "incident_incident_id",
				Unique:  false,
				Columns: []*schema.Column{IncidentsColumns[6]},
			},
			{
				Name:    "incident_name",
				Unique:  false,
				Columns: []*schema.Column{IncidentsColumns[7]},
			},
		},
	}
	// StatusColumns holds the columns for the "status" table.
	StatusColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "value", Type: field.TypeString},
		{Name: "alert_stati", Type: field.TypeInt, Nullable: true},
		{Name: "incident_stati", Type: field.TypeInt, Nullable: true},
	}
	// StatusTable holds the schema information for the "status" table.
	StatusTable = &schema.Table{
		Name:       "status",
		Columns:    StatusColumns,
		PrimaryKey: []*schema.Column{StatusColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "status_alerts_Stati",
				Columns:    []*schema.Column{StatusColumns[3]},
				RefColumns: []*schema.Column{AlertsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "status_incidents_Stati",
				Columns:    []*schema.Column{StatusColumns[4]},
				RefColumns: []*schema.Column{IncidentsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AlertsTable,
		CountersTable,
		FailuresTable,
		FilesTable,
		IncidentsTable,
		StatusTable,
	}
)

func init() {
	CountersTable.ForeignKeys[0].RefTable = AlertsTable
	CountersTable.ForeignKeys[1].RefTable = IncidentsTable
	FailuresTable.ForeignKeys[0].RefTable = AlertsTable
	FailuresTable.ForeignKeys[1].RefTable = IncidentsTable
	FilesTable.ForeignKeys[0].RefTable = AlertsTable
	FilesTable.ForeignKeys[1].RefTable = IncidentsTable
	StatusTable.ForeignKeys[0].RefTable = AlertsTable
	StatusTable.ForeignKeys[1].RefTable = IncidentsTable
}

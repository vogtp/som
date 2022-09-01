package webstatus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// MigrateIncidents TODO remove
func (s *WebStatus) MigrateIncidents() {
	a := s.DB()
	ctx := context.Background()
	files, err := s.getIncidentDetailFiles(s.getIncidentRoot(), "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %v incident files\n", len(files))
	for _, f := range files {
		if len(f.IncidentInfo.ByteState) < 1 {
			fmt.Println("NO BYTE STATE in files")
		}
		if err := a.SaveIncident(ctx, f.IncidentInfo.IncidentMsg); err != nil {
			panic(err)
		}

	}
}

// Incidents TODO remove
func (s *WebStatus) Incidents() {
	a := s.DB()

	incidents, err := a.GetIncident(context.Background(), "")
	if err != nil {
		panic(err)
	}
	for _, r := range incidents {
		fmt.Printf("%s %-20s %10v %v %v %v\n", r.IncidentID, r.Name, r.Level(), r.Start, r.End, r.Error)
	}
	fmt.Printf("Total incidents: %v\n", len(incidents))
}

// Files TODO remove
func (s *WebStatus) Files(pid uuid.UUID) {
	a := s.DB()

	files, err := a.GetFiles(context.Background(), pid)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Printf("%s %-20s %10v %v \n", f.Name, f.Type.Ext, f.Size, f.Type.MimeType)
	}
	fmt.Printf("Total files: %v\n", len(files))
}

// IncidentsSummary TODO remove
func (s *WebStatus) IncidentsSummary() {
	a := s.DB()

	summary, err := a.GetIncidentSummary(context.Background(), "")
	if err != nil {
		panic(err)
	}
	for _, r := range summary {
		fmt.Printf("%s %-20s %10v %v %v %v\n", r.IncidentID, r.Name, r.Level(), r.Start, r.End, r.Error)
	}
	fmt.Printf("Total summaries: %v\n", len(summary))
}

package webstatus

import (
	"fmt"

	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

// MigrateIncidents ..
func (s *WebStatus) MigrateIncidents() {
	a := db.Access{}
	files, err := s.getIncidentDetailFiles(s.getIncidentRoot(), "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %v incident files\n", len(files))
	for _, f := range files {
		if len(f.IncidentInfo.ByteState) < 1 {
			fmt.Println("NO BYTE STATE in files")
		}
		if err := a.SaveIncident(f.IncidentInfo.IncidentMsg); err != nil {
			panic(err)
		}

	}
}

func (s *WebStatus) Query() {
	// s.Incidents()
	s.Summay()
}

func (s *WebStatus) Incidents() {
	a := db.Access{}

	incidents, err := a.GetIncident("")
	if err != nil {
		panic(err)
	}
	for _, r := range incidents {
		fmt.Printf("%s %-20s %10v %v %v %v\n", r.IncidentID, r.Name, r.Level(), r.Start, r.End, r.Error)
	}
	fmt.Printf("Total incidents: %v\n", len(incidents))
}

func (s *WebStatus) Summay() {
	a := db.Access{}

	summary, err := a.GetIncidentSummary("")
	if err != nil {
		panic(err)
	}
	for _, r := range summary {
		fmt.Printf("%s %-20s %10v %v %v %v\n", r.IncidentID, r.Name, r.Level(), r.Start, r.End, r.Error)
	}
	fmt.Printf("Total summaries: %v\n", len(summary))
}

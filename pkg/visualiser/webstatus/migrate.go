package webstatus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/cfg"
)

// MigrateIncidents TODO remove
func (s *WebStatus) MigrateIncidents() {
	a := s.DB()
	ent := s.Ent()
	//defer ent.Close()
	//ctx := context.Background()
	_, err := s.getIncidentDetailFiles(a, ent, s.getIncidentRoot(), "")
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Migrated %v incident files\n", len(files))
	// for _, f := range files {
	// 	if len(f.IncidentInfo.ByteState) < 1 {
	// 		fmt.Println("NO BYTE STATE in files")
	// 	}
	// 	if err := a.SaveIncident(ctx, f.IncidentInfo.IncidentMsg); err != nil {
	// 		panic(err)
	// 	}

	// }
}

// MigrateAlerts TODO remove
func (s *WebStatus) MigrateAlerts() {
	a := s.DB()
	ent := s.Ent()
	//defer ent.Close()
	_, err := s.getAlertFiles(a, ent, s.getAlertRoot(), "")
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Migrated %v alert files\n", len(files))
	// errCnt := 0
	// for _, f := range files {
	// 	alert, err := s.getAlert(f.Path)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if err := a.SaveAlert(ctx, alert); err != nil {
	// 		hcl.Errorf("Cannot save alert %s: %v", alert.ID.String(), err)
	// 		errCnt++
	// 	}

	// }
	// hcl.Infof("Got %v/%v errors", errCnt, len(files))
}

// Alerts TODO remove
func (s *WebStatus) Alerts() {
	a := s.DB()

	alerts, err := a.GetAlert(context.Background(), "")
	if err != nil {
		panic(err)
	}
	for _, r := range alerts {
		fmt.Printf("%s %-20s %10v %v \n", r.Time.Format(cfg.TimeFormatString), r.Name, r.Level(), r.Error)
	}
	fmt.Printf("Total alerts: %v\n", len(alerts))
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

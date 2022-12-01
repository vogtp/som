package webstatus

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/status"
)

var onceDB sync.Once

func (s *WebStatus) cleanup() {
	cleanupIntervall := 6 * time.Hour
	cleanupTimeout := time.Hour
	onceDB.Do(func() {
		ticker := time.NewTicker(6 * time.Hour)
		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
				s.thinOutIncidents(ctx)
				s.cleanupIncidents(ctx)
				cancel()
				// cleanup every 6h
				time.Sleep(cleanupIntervall)
				<-ticker.C
			}
		}()
	})
}

func (s *WebStatus) thinOutIncidents(ctx context.Context) {
	if err := s.dbAccess.ThinOutIncidents(ctx); err != nil {
		s.hcl.Warnf("thining out incidents failed: %v", err)
	}
}
func (s *WebStatus) cleanupIncidents(ctx context.Context) {
	ent := s.Ent()
	all, err := ent.IncidentSummary.Query().All(ctx)
	if err != nil {
		s.hcl.Warnf("Cannot close stale incidents: %v", err)
		return
	}
	for _, inci := range all {
		if !inci.End.IsZero() {
			continue
		}
		lvl := status.Unknown
		if szGrp := s.data.Status.Get(inci.Name); szGrp != nil {
			lvl = szGrp.Level()
		}
		if lvl > status.OK {
			// not cleaning up since status is not OK or UNKNOWN
			continue
		}
		s.hcl.Infof("Closing incident: %T %v -> %s\n", inci, inci.Name, lvl)
		closer := ent.Incident.IncidentClient.Create()
		closer.SetUUID(uuid.New())
		closer.SetIncidentID(inci.IncidentID)
		closer.SetName(inci.Name)
		// fix stupid time formating
		now, _ := time.Parse(cfg.TimeFormatString, time.Now().Format(cfg.TimeFormatString))
		closer.SetTime(now)
		closer.SetStart(inci.Start.Time())
		closer.SetEnd(now)
		closer.SetIntLevel(int(lvl))
		closer.SetError("Autoclosed")
		closer.SetUsername("Cleanup Job")
		closer.SetRegion("")
		closer.SetProbeHost("")
		closer.SetProbeOS("")
		closer.SetState([]byte(""))
		if err := closer.Exec(ctx); err != nil {
			s.hcl.Warnf("Cannot save incident %s %v: %v", inci.Name, inci.IncidentID, err)
		}
	}
}

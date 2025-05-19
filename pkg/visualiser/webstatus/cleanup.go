package webstatus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
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
				if s.dbAccess != nil {
					ctx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
					s.thinOutIncidents(ctx)
					s.cleanupIncidents(ctx)
					cancel()
				}
				// cleanup every 6h
				time.Sleep(cleanupIntervall)
				<-ticker.C
			}
		}()
	})
}

func (s *WebStatus) thinOutIncidents(ctx context.Context) {
	if err := s.dbAccess.ThinOutIncidents(ctx); err != nil {
		s.log.Warn("thining out incidents failed", log.Error, err)
	}
}
func (s *WebStatus) cleanupIncidents(ctx context.Context) {
	incidentSummary := s.Ent().IncidentSummary
	all, err := incidentSummary.Query().All(ctx)
	if err != nil {
		s.log.Warn("Cannot close stale incidents", log.Error, err)
		return
	}

	autocloseDuration := 2 * viper.GetDuration(cfg.AlertIncidentCorrelationReopenTime)
	s.log.Info("Cleaning up stale incidents", "autoclose", autocloseDuration)
	for _, is := range all {
		if !is.End.IsZero() {
			continue
		}
		lvl := status.Unknown
		if szGrp := s.data.Status.Get(is.Name); szGrp != nil {
			lvl = szGrp.Level()
		}
		if lvl > status.OK && time.Since(is.LastUpdate.Time()) < autocloseDuration {
			// not cleaning up since status is not OK or UNKNOWN
			continue
		}
		s.log.Info("Closing incident", log.Szenario, is.Name, "level", lvl)
		err := incidentSummary.CloseIncident(ctx, is, "Cleanup Job", lvl, fmt.Sprintf("%s (Stale incident: autoclosed)", is.Error))
		if err != nil {
			s.log.Warn("Cannot save incident %s %v: %v", log.Szenario, is.Name, "incident_id", is.IncidentID, log.Error, err)
		}
	}
}

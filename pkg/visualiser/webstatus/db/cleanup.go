package db

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

const (
	maxIncidentsPerSummary    = 200
	maxIncidentsPerSummary30d = 50
)

var cleanupMutex sync.Mutex

// ThinOutIncidents removes multiple incident entries
func (c *Client) ThinOutIncidents(ctx context.Context) error {
	c.hcl.Info("Starting thin out")
	defer c.hcl.Info("Finish thin out")
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()
	incidents, err := c.IncidentSummary.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("cannot query incident summaries: %w", err)
	}
	for _, s := range incidents {
		if s.Total > maxIncidentsPerSummary {
			if err := c.thinoutIncident(ctx, s, maxIncidentsPerSummary); err != nil {
				return fmt.Errorf("cannot thin out indicent %v: %w", s.IncidentID, err)
			}
			continue
		}
		if !s.End.t.IsZero() && time.Since(s.End.t) < 30*24*time.Hour {
			continue
		}
		if s.Total > maxIncidentsPerSummary30d {
			if err := c.thinoutIncident(ctx, s, maxIncidentsPerSummary30d); err != nil {
				return fmt.Errorf("cannot thin out indicent %v: %w", s.IncidentID, err)
			}
			continue
		}
	}
	return nil
}

func (c *Client) thinoutIncident(ctx context.Context, incSum *IncidentSummary, maxIncidents int) error {
	intervall := int(math.Ceil(float64(incSum.Total) / float64(maxIncidents)))
	c.hcl.Info("Thining out", "incident", incSum.Name, "total_entries", incSum.Total, "intervall", intervall, "incident_id", incSum.IncidentID)
	incidents, err := c.Incident.Query().Where(incident.IncidentIDEQ(incSum.IncidentID)).Order(ent.Asc(incident.FieldTime)).All(ctx)
	if err != nil {
		return fmt.Errorf("cannot query incidents of %v: %w", incSum.IncidentID, err)
	}
	i := 0
	thisFailure := ""
	for _, inc := range incidents {
		fail, err := inc.QueryFailures().Order(ent.Desc(failure.FieldIdx)).First(ctx)
		if err != nil {
			c.hcl.Debug("cannot get failures","error", err)
			continue
		}
		lastFailure := thisFailure
		thisFailure = fail.Error
		if lastFailure != thisFailure {
			continue
		}
		i--
		if i < 1 {
			i = intervall
			continue
		}
		if _, err := c.Incident.Delete().Where(incident.ID(inc.ID)).Exec(ctx); err != nil {
			return fmt.Errorf("error deleting incident %v: %w", inc.ID, err)
		}
	}

	return nil
}

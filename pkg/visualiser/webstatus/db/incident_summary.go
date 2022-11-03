package db

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/failure"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/predicate"
)

const (
	maxIncidentsPerSummary    = 200
	maxIncidentsPerSummary30d = 50
)

// IncidentSummary is the summary of different incidents entries
type IncidentSummary struct {
	IncidentID uuid.UUID  `json:"incident_id"`
	Name       string     `json:"name"`
	Total      int        `json:"total"`
	IntLevel   int        `json:"level"`
	Start      MinMaxTime `json:"start"`
	End        MinMaxTime `json:"end"`
	Error      string     `json:"error"`
}

// IncidentSummaryQuery allows quering incidents
type IncidentSummaryQuery struct {
	client *Client
	*ent.IncidentQuery
}

// Query returns a list of incidents
func (isq *IncidentSummaryQuery) Query() *IncidentSummaryQuery {
	q := isq.client.Incident.Query()
	isq.IncidentQuery = q
	return isq
}

func (isq *IncidentSummaryQuery) Where(s ...predicate.Incident) *IncidentSummaryQuery {
	isq.IncidentQuery.Where(s...)
	return isq
}

func (isq *IncidentSummaryQuery) First(ctx context.Context) (*IncidentSummary, error) {
	isq.IncidentQuery.Limit(1)
	all, err := isq.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, &ent.NotFoundError{}
	}
	return all[0], nil
}

// All returns all incudent summaries
func (isq *IncidentSummaryQuery) All(ctx context.Context) ([]*IncidentSummary, error) {
	var summary []*IncidentSummary
	err := isq.groupAndAggregate().Scan(ctx, &summary)
	// last event is the OK so remove it
	for i, s := range summary {
		summary[i].Total = s.Total - 1
	}
	return summary, err
}

func (isq *IncidentSummaryQuery) groupAndAggregate() *ent.IncidentGroupBy {
	return isq.IncidentQuery.Select(incident.FieldIncidentID, incident.FieldName).
		GroupBy(incident.FieldIncidentID, incident.FieldName).
		Aggregate(
			ent.As(ent.Count(), "Total"),
			ent.As(ent.Max(incident.FieldIntLevel), "Level"),
			ent.As(ent.Max(incident.FieldEnd), "End"),
			ent.As(ent.Min(incident.FieldStart), "Start"),
			ent.As(ent.Max(incident.FieldError), "Error"),
		)
}

// Level convinience method that calls status Level
func (is IncidentSummary) Level() status.Level {
	return status.Level(is.IntLevel)
}

// ThinOut removes doublicte entries
func (isq *IncidentSummaryQuery) ThinOut(ctx context.Context) error {
	incidents, err := isq.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("cannot query incident summaries: %w", err)
	}
	for _, s := range incidents {
		if s.Total > maxIncidentsPerSummary {
			if err := isq.thinoutIncident(ctx, s, maxIncidentsPerSummary); err != nil {
				return fmt.Errorf("cannot thin out indicent %v: %w", s.IncidentID, err)
			}
			continue
		}
		if ! s.End.t.IsZero() && time.Since(s.End.t) < 30*24*time.Hour {
			continue
		}
		if s.Total > maxIncidentsPerSummary30d {
			if err := isq.thinoutIncident(ctx, s, maxIncidentsPerSummary30d); err != nil {
				return fmt.Errorf("cannot thin out indicent %v: %w", s.IncidentID, err)
			}
			continue
		}
	}
	return nil
}

func (isq *IncidentSummaryQuery) thinoutIncident(ctx context.Context, incSum *IncidentSummary, maxIncidents int) error {
	intervall := int(math.Ceil(float64(incSum.Total) / float64(maxIncidents)))
	isq.client.hcl.Infof("Thining out: %s %v (every %v) %v", incSum.Name, incSum.Total, intervall, incSum.IncidentID)
	incidents, err := isq.client.Incident.Query().Where(incident.IncidentIDEQ(incSum.IncidentID)).Order(ent.Asc(incident.FieldTime)).All(ctx)
	if err != nil {
		return fmt.Errorf("cannot query incidents of %v: %w", incSum.IncidentID, err)
	}
	i := intervall
	thisFailure := ""
	notFound := ent.NotFoundError{}
	for _, inc := range incidents {
		fail, err := inc.QueryFailures().Order(ent.Desc(failure.FieldIdx)).First(ctx)
		if err != nil {
			if !errors.Is(err, &notFound) {
				isq.client.hcl.Warnf("cannot get failures: %v", err)
			}
			continue
		}
		lastFailure := thisFailure
		thisFailure = fail.Error
		if lastFailure != thisFailure {
			isq.client.hcl.Tracef("%q != %q", thisFailure, lastFailure)
			continue
		}
		i--
		if i < 1 {
			i = intervall
			continue
		}
		if _, err := isq.client.Incident.Delete().Where(incident.ID(inc.ID)).Exec(ctx); err != nil {
			return fmt.Errorf("error deleting incident %v: %w", inc.ID, err)
		}
	}

	return nil
}

func (isq *IncidentSummaryQuery) reportChilds(ctx context.Context, inc *ent.Incident) {
	fmt.Printf("Incident %v children:\n", inc.ID)
	file, err := inc.QueryFiles().All(ctx)
	if err != nil {
		isq.client.hcl.Errorf("report child: %v", err)
	}
	fmt.Printf("files %v\n", len(file))
	cntrs, err := inc.QueryCounters().All(ctx)
	if err != nil {
		isq.client.hcl.Errorf("report child: %v", err)
	}
	fmt.Printf("counters %v\n", len(cntrs))
	fails, err := inc.QueryFailures().All(ctx)
	if err != nil {
		isq.client.hcl.Errorf("report child: %v", err)
	}
	fmt.Printf("fails %v\n", len(fails))
}

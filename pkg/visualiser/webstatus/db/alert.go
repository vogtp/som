package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/alert"
)

// AlertClient is a wrapper enhaning the ent client
type AlertClient struct {
	*ent.AlertClient
	client *Client
}

// Szenarios returns a list of szenario names
func (ac *AlertClient) Szenarios(ctx context.Context) ([]string, error) {
	return ac.client.Alert.Query().Select(alert.FieldName).GroupBy(alert.FieldName).Strings(ctx)
}

// Save an alert msg to ent
func (ac *AlertClient) Save(ctx context.Context, msg *msg.AlertMsg) error {

	i := ac.Create()

	i.SetUUID(msg.ID)
	if incID, err := uuid.Parse(msg.IncidentID); err == nil {
		i.SetIncidentID(incID)
	}
	i.SetName(msg.Name)
	i.SetTime(msg.Time)
	i.SetUsername(msg.Username)
	i.SetRegion(msg.Region)
	i.SetProbeOS(msg.ProbeOS)
	i.SetProbeHost(msg.ProbeHost)
	if msg.Err() != nil {
		i.SetError(msg.Err().Error())
	}
	i.SetIntLevel(int(status.Unknown.FromString(msg.Level)))

	if errs, err := ac.client.getErrors(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFailures(errs...)
	} else {
		ac.client.log.Warn("Getting errors", log.Error, err)
	}
	if stati, err := ac.client.getStati(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddStati(stati...)
	} else {
		ac.client.log.Warn("Getting stari", log.Error, err)
	}
	if cntrs, err := ac.client.getCounter(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddCounters(cntrs...)
	} else {
		ac.client.log.Warn("Getting counters", log.Error, err)
	}
	if fils, err := ac.client.getFiles(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFiles(fils...)
	} else {
		ac.client.log.Warn("Getting files", log.Error, err)
	}

	return i.Exec(ctx)
}

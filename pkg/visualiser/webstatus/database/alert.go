package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent"
)

// AlertClient is a wrapper enhaning the ent client
type AlertClient struct {
	*ent.AlertClient
	client *Client
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
	i.SetLevel(int(status.Unknown.FromString(msg.Level)))

	if errs, err := ac.client.getErrors(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFailures(errs...)
	} else {
		ac.client.hcl.Warnf("Getting errors: %v", err)
	}
	if stati, err := ac.client.getStati(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddStati(stati...)
	} else {
		ac.client.hcl.Warnf("Getting stari: %v", err)
	}
	if cntrs, err := ac.client.getCounter(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddCounters(cntrs...)
	} else {
		ac.client.hcl.Warnf("Getting counters: %v", err)
	}
	if fils, err := ac.client.getFiles(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFiles(fils...)
	} else {
		ac.client.hcl.Warnf("Getting files: %v", err)
	}

	return i.Exec(ctx)
}

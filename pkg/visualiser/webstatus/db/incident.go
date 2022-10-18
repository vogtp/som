package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

// IncidentClient is a wrapper enhaning the ent client
type IncidentClient struct {
	*ent.IncidentClient
	client *Client
}

// Szenarios returns a list of szenario names
func (ic *IncidentClient) Szenarios(ctx context.Context) ([]string, error) {
	return ic.client.Incident.Query().Select(incident.FieldName).GroupBy(incident.FieldName).Strings(ctx)
}

// Save save an incident msg to ent
func (ic *IncidentClient) Save(ctx context.Context, msg *msg.IncidentMsg) error {

	i := ic.IncidentClient.Create()

	i.SetUUID(msg.ID)
	if incID, err := uuid.Parse(msg.IncidentID); err == nil {
		i.SetIncidentID(incID)
	}else{
		ic.client.hcl.Errorf("Cannot parse incident ID: %v",err)
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
	if !msg.Start.IsZero() {
		i.SetStart(msg.Start)
	}
	if !msg.End.IsZero() {
		i.SetEnd(msg.End)
	}
	i.SetIntLevel(msg.IntLevel)
	i.SetState(msg.ByteState)

	if errs, err := ic.client.getErrors(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFailures(errs...)
	} else {
		ic.client.hcl.Warnf("Getting errors: %v", err)
	}
	if stati, err := ic.client.getStati(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddStati(stati...)
	} else {
		ic.client.hcl.Warnf("Getting stari: %v", err)
	}
	if cntrs, err := ic.client.getCounter(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddCounters(cntrs...)
	} else {
		ic.client.hcl.Warnf("Getting counters: %v", err)
	}
	if fils, err := ic.client.getFiles(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFiles(fils...)
	} else {
		ic.client.hcl.Warnf("Getting files: %v", err)
	}

	return i.Exec(ctx)
}

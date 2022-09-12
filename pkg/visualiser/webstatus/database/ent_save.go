package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/ent"
)

// SaveIncident save an incident msg to ent
func (a *Access) SaveIncident(ctx context.Context, msg *msg.IncidentMsg) error {

	i := a.client.Incident.Create()

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
	i.SetStart(msg.Start)
	i.SetEnd(msg.End)
	i.SetLevel(msg.IntLevel)
	i.SetState(msg.ByteState)

	if errs, err := a.getErrors(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFailures(errs...)
	} else {
		a.hcl.Warnf("Getting errors: %v", err)
	}
	if stati, err := a.getStati(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddStati(stati...)
	} else {
		a.hcl.Warnf("Getting stari: %v", err)
	}
	if cntrs, err := a.getCounter(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddCounters(cntrs...)
	} else {
		a.hcl.Warnf("Getting counters: %v", err)
	}
	if fils, err := a.getFiles(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFiles(fils...)
	} else {
		a.hcl.Warnf("Getting files: %v", err)
	}

	return i.Exec(ctx)
}

// SaveAlert save an alert msg to ent
func (a *Access) SaveAlert(ctx context.Context, msg *msg.AlertMsg) error {

	i := a.client.Alert.Create()

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

	if errs, err := a.getErrors(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFailures(errs...)
	} else {
		a.hcl.Warnf("Getting errors: %v", err)
	}
	if stati, err := a.getStati(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddStati(stati...)
	} else {
		a.hcl.Warnf("Getting stari: %v", err)
	}
	if cntrs, err := a.getCounter(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddCounters(cntrs...)
	} else {
		a.hcl.Warnf("Getting counters: %v", err)
	}
	if fils, err := a.getFiles(ctx, msg.SzenarioEvtMsg); err == nil {
		i.AddFiles(fils...)
	} else {
		a.hcl.Warnf("Getting files: %v", err)
	}

	return i.Exec(ctx)
}

func (a *Access) getErrors(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.Failure, error) {
	var reterr error
	i := 0
	errs := make([]*ent.Failure, len(msg.Errors))
	for idx, e := range msg.Errs() {
		t, err := a.client.Failure.Create().SetIdx(idx).SetError(e).Save(ctx)
		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				err = fmt.Errorf("%v %w", reterr, err)
			}
		}
		errs[i] = t
		i++
	}
	return errs, reterr
}

func (a *Access) getStati(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.Status, error) {
	var reterr error
	i := 0
	stati := make([]*ent.Status, len(msg.Stati))
	for k, v := range msg.Stati {
		t, err := a.client.Status.Create().SetName(k).SetValue(v).Save(ctx)
		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				err = fmt.Errorf("%v %w", reterr, err)
			}
		}
		stati[i] = t
		i++
	}
	return stati, reterr
}

func (a *Access) getCounter(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.Counter, error) {
	var reterr error
	i := 0
	cntr := make([]*ent.Counter, len(msg.Counters))
	for k, v := range msg.Counters {
		t, err := a.client.Counter.Create().SetName(k).SetValue(fmt.Sprintf("%v", v)).Save(ctx)
		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				err = fmt.Errorf("%v %w", reterr, err)
			}
		}
		cntr[i] = t
		i++
	}
	return cntr, reterr
}

func (a *Access) getFiles(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.File, error) {
	var reterr error
	fils := make([]*ent.File, len(msg.Files))
	for i, f := range msg.Files {
		f.CalculateID()
		t, err := a.client.File.Create().
			SetUUID(f.ID).SetName(f.Name).SetType(f.Type.MimeType).SetExt(f.Type.Ext).
			SetPayload(f.Payload).SetSize(f.Size).
			Save(ctx)

		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				err = fmt.Errorf("%v %w", reterr, err)
			}
		}
		fils[i] = t
		i++
	}
	return fils, reterr
}

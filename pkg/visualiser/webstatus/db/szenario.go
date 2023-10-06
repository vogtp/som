package db

import (
	"context"
	"fmt"

	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
)

func (client *Client) getErrors(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.Failure, error) {
	var reterr error
	i := 0
	errs := make([]*ent.Failure, len(msg.Errors))
	for idx, e := range msg.Errs() {
		t, err := client.Failure.Create().SetIdx(idx).SetError(e).Save(ctx)
		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				reterr = fmt.Errorf("%v %w", reterr, err)
			}
		}
		errs[i] = t
		i++
	}
	return errs, reterr
}

func (client *Client) getStati(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.Status, error) {
	var reterr error
	i := 0
	stati := make([]*ent.Status, len(msg.Stati))
	for k, v := range msg.Stati {
		t, err := client.Status.Create().SetName(k).SetValue(v).Save(ctx)
		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				reterr = fmt.Errorf("%v %w", reterr, err)
			}
		}
		stati[i] = t
		i++
	}
	return stati, reterr
}

func (client *Client) getCounter(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.Counter, error) {
	var reterr error
	i := 0
	cntr := make([]*ent.Counter, len(msg.Counters))
	for k, v := range msg.Counters {
		t, err := client.Counter.Create().SetName(k).SetValue(v).Save(ctx)
		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				reterr = fmt.Errorf("%v %w", reterr, err)
			}
		}
		cntr[i] = t
		i++
	}
	return cntr, reterr
}

func (client *Client) getFiles(ctx context.Context, msg *msg.SzenarioEvtMsg) ([]*ent.File, error) {
	var reterr error
	fils := make([]*ent.File, len(msg.Files))
	for i, f := range msg.Files {
		f := f
		f.CalculateID()
		t, err := client.File.Create().
			SetUUID(f.ID).SetName(f.Name).SetType(f.Type.MimeType).SetExt(f.Type.Ext).
			SetPayload(f.Payload).SetSize(f.Size).
			Save(ctx)

		if err != nil {
			if reterr == nil {
				reterr = err
			} else {
				reterr = fmt.Errorf("%v %w", reterr, err)
			}
		}
		fils[i] = t
	}
	return fils, reterr
}

package db

import (
	"context"
	"fmt"

	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

// SaveAlert saves a alert to DB
func (a *Access) SaveAlert(ctx context.Context, msg *msg.AlertMsg) error {
	db := a.getDb()
	var reterr error
	if err := a.SaveCounters(ctx, msg.SzenarioEvtMsg); err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}
	if err := a.SaveStati(ctx, msg.SzenarioEvtMsg); err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}
	if err := a.SaveErrors(ctx, msg.SzenarioEvtMsg); err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}
	if err := a.SaveFiles(ctx, msg.SzenarioEvtMsg); err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}
	model := &AlertModel{
		IntLevel:      int(status.Unknown.FromString(msg.Level)),
		SzenarioModel: a.SzenarioModelFromMsg(msg.SzenarioEvtMsg),
	}
	if err := db.WithContext(ctx).Save(model).Error; err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}
	return reterr
}

// AlertModel the DB model of a alert (use msg?)
type AlertModel struct {
	SzenarioModel
	IntLevel int `json:"Level" gorm:"column:Level"`
}

// Level convinience method that calls status Level
func (im AlertModel) Level() status.Level {
	return status.Level(im.IntLevel)
}

// AlertSzenarios lists all szenarios that have alerts
func (a *Access) AlertSzenarios(ctx context.Context) []string {
	db := a.getDb()
	result := make([]string, 0)
	db.Model(&AlertModel{}).Distinct("name").Order("name COLLATE NOCASE").WithContext(ctx).Find(&result)
	return result
}

// GetAlert returns a alert list by id (uuid)
func (a *Access) GetAlert(ctx context.Context, id string) ([]AlertModel, error) {
	db := a.getDb()
	result := make([]AlertModel, 0)
	search := db.Model(&AlertModel{}).Order("time")
	if len(id) > 0 {
		search = search.Where("id = ?", id)
	}
	err := search.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load alert: %w", err)
	}
	return result, err
}

// GetAlertBySzenario returns alerts list by szenario name
func (a *Access) GetAlertBySzenario(ctx context.Context, sz string) ([]AlertModel, error) {
	db := a.getDb()
	result := make([]AlertModel, 0)
	search := db.Model(&AlertModel{}).Order("time desc")
	if len(sz) > 0 && sz != "all" {
		search = search.Where("name like ?", sz)
	}
	err := search.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load alert: %w", err)
	}
	return result, err
}

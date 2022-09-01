package db

import (
	"context"
	"fmt"
	"time"

	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

// SaveIncident saves a incident to DB
func (a *Access) SaveIncident(ctx context.Context, msg *msg.IncidentMsg) error {
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
	model := IncidentModel{
		Start:         msg.Start,
		End:           msg.End,
		IntLevel:      msg.IntLevel,
		ByteState:     msg.ByteState,
		SzenarioModel: a.SzenarioModelFromMsg(msg.SzenarioEvtMsg),
	}
	if err := db.WithContext(ctx).Save(model).Error; err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}
	if err := db.WithContext(ctx).Model(&model).Where("id = ?", model.ID).Update("end", model.End).Error; err != nil {
		a.hcl.Warnf("Cannot update incident end times: %v", err)
	}
	return reterr
}

// IncidentModel the DB model of a incident (use msg?)
type IncidentModel struct {
	SzenarioModel
	Start     time.Time `gorm:"index"`
	End       time.Time `gorm:"index"`
	IntLevel  int       `json:"Level" gorm:"column:Level"`
	ByteState []byte    `json:"State" gorm:"column:State"`
}

// Level convinience method that calls status Level
func (im IncidentModel) Level() status.Level {
	return status.Level(im.IntLevel)
}

// IncidentSummary db model for the incident list
type IncidentSummary struct {
	IncidentID string
	Name       string
	Start      MinMaxTime
	End        MinMaxTime
	Total      int
	IntLevel   int
	Error      string
	DetailLink string
}

// Level convinience method that calls status Level
func (il IncidentSummary) Level() status.Level {
	return status.Level(il.IntLevel)
}

// IncidentSzenarios lists all szenarios that have incidents
func (a *Access) IncidentSzenarios(ctx context.Context) []string {
	db := a.getDb()
	result := make([]string, 0)
	db.Model(&IncidentModel{}).Distinct("name").Order("name COLLATE NOCASE").WithContext(ctx).Find(&result)
	return result
}

// GetIncident returns a incident list by id (uuid)
func (a *Access) GetIncident(ctx context.Context, id string) ([]IncidentModel, error) {
	db := a.getDb()
	result := make([]IncidentModel, 0)
	search := db.Model(&IncidentModel{}).Order("time")
	if len(id) > 0 {
		search = search.Where("incident_id = ?", id)
	}
	err := search.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load incident: %w", err)
	}
	return result, err
}

// GetIncidentSummary returns a incident summary list by szeanrio name
func (a *Access) GetIncidentSummary(ctx context.Context, szName string) ([]IncidentSummary, error) {
	db := a.getDb()
	result := make([]IncidentSummary, 0)
	search := db.Model(&IncidentModel{}).Select("incident_id, name, count(*) as Total, MAX(Level) as IntLevel, MAX(end) as End, MIN(start) as Start, MAX(Error) as Error")
	search = search.Group("incident_id").Order("Start desc")
	if len(szName) > 1 && szName != "all" {
		search = search.Where("name like ?", szName)
	}
	err := search.WithContext(ctx).Find(&result).Error
	return result, err
}

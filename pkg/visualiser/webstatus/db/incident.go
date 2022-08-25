package db

import (
	"fmt"
	"time"

	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

func (a *Access) SaveIncident(msg *msg.IncidentMsg) error {
	db := a.getDb()
	var reterr error
	if err := a.SaveStati(msg.SzenarioEvtMsg); err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}
	if err := a.SaveErrors(msg.SzenarioEvtMsg); err != nil {
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
	if err := db.Save(model).Error; err != nil {
		if reterr == nil {
			reterr = err
		} else {
			err = fmt.Errorf("%v %w", reterr, err)
		}
	}

	return reterr
}

type IncidentModel struct {
	SzenarioModel
	Start     time.Time `gorm:"index"`
	End       time.Time `gorm:"index"`
	IntLevel  int       `json:"Level" gorm:"column:Level"`
	ByteState []byte    `json:"State" gorm:"column:State"`
}

func (im IncidentModel) Level() status.Level {
	return status.Level(im.IntLevel)
}

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

func (il IncidentSummary) Level() status.Level {
	return status.Level(il.IntLevel)
}

func (a *Access) IncidentSzenarios() []string {
	db := a.getDb()
	result := make([]string, 0)
	db.Model(&IncidentModel{}).Distinct("name").Find(&result)
	return result
}

func (a *Access) GetIncidents() ([]IncidentModel, error) {
	return a.GetIncident("")
}

func (a *Access) GetIncident(id string) ([]IncidentModel, error) {
	db := a.getDb()
	result := make([]IncidentModel, 0)
	search := db.Model(&IncidentModel{}).Order("time")
	if len(id) > 0 {
		search = search.Where("incident_id = ?", id)
	}
	err := search.Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load incident: %w", err)
	}
	return result, err
}

func (a *Access) GetIncidentSummary(szName string) ([]IncidentSummary, error) {
	db := a.getDb()
	result := make([]IncidentSummary, 0)
	search := db.Model(&IncidentModel{}).Select("incident_id, name, count(*) as Total, MAX(Level) as IntLevel, MAX(time) as End, MIN(time) as Start").Group("incident_id").Order("Start")
	if len(szName) > 1 && szName != "all" {
		search = search.Where("name like ?", szName)
	}
	err := search.Find(&result).Error

	for i, r := range result {
		incs, err := a.GetIncident(r.IncidentID)
		if err != nil {
			a.hcl.Infof("incident not found: %v", err)
			continue
		}
		for _, inc := range incs {
			if len(inc.Error) < 1 {
				continue
			}
			result[i].Error = inc.Error
		}
	}
	return result, err
}

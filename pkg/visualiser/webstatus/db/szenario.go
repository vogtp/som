package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/msg"
)

type SzenarioModel struct {
	ID         uuid.UUID `json:"ID"  gorm:"primaryKey;type:uuid"`
	IncidentID string    `json:"Incident" gorm:"index"`
	Name       string    `json:"Name" gorm:"index"`
	Time       time.Time `json:"Time" gorm:"index"`
	Username   string    `json:"Username"`
	Region     string    `json:"Region"`

	// FIXME handle:
	// Files     []msg.FileMsgItem `json:"Files" gorm:"foreignKey:MsgID;references:ID"`
	ProbeOS   string `json:"OS"`
	ProbeHost string `json:"Host"`
	Error     string
}

type stati struct {
	ParentID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name     string
	Value    string
}

type ErrorModel struct {
	ParentID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Idx      int
	Error    string
}

func (a *Access) GetErrors(id uuid.UUID) ([]ErrorModel, error) {
	db := a.getDb()
	result := make([]ErrorModel, 0)
	search := db.Model(&ErrorModel{}).Order("idx")
	if len(id) > 0 {
		search = search.Where("parent_id = ?", id)
	}
	err := search.Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load error: %w", err)
	}
	return result, err
}

func (a Access) SzenarioModelFromMsg(msg *msg.SzenarioEvtMsg) SzenarioModel {
	sm := SzenarioModel{
		ID:         msg.ID,
		IncidentID: msg.IncidentID,
		Name:       msg.Name,
		Time:       msg.Time,
		Username:   msg.Username,
		Region:     msg.Region,
		ProbeOS:    msg.ProbeOS,
		ProbeHost:  msg.ProbeHost,
	}
	if msg.Err() != nil {
		sm.Error = msg.Err().Error()
	}
	return sm
}

func (a *Access) SaveErrors(msg *msg.SzenarioEvtMsg) error {
	db := a.getDb()
	var reterr error
	for i, e := range msg.Errors {
		if err := db.Save(ErrorModel{
			ParentID: msg.ID,
			Idx:      i,
			Error:    e,
		}).Error; err != nil {
			if reterr == nil {
				reterr = err
			} else {
				err = fmt.Errorf("%v %w", reterr, err)
			}
		}
	}
	return reterr
}

func (a *Access) SaveStati(msg *msg.SzenarioEvtMsg) error {
	db := a.getDb()
	var reterr error
	for k, v := range msg.Stati {
		if err := db.Save(stati{
			ParentID: msg.ID,
			Name:     k,
			Value:    v,
		}).Error; err != nil {
			if reterr == nil {
				reterr = err
			} else {
				err = fmt.Errorf("%v %w", reterr, err)
			}
		}
	}
	return reterr
}

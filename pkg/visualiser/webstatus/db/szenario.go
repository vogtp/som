package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/msg"
)

// SzenarioModel model for szenarios
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

type statiModel struct {
	ParentID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name     string
	Value    string
}

type counterModel struct {
	ParentID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name     string
	Value    string
}

// ErrorModel model for errors
type ErrorModel struct {
	ParentID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Idx      int
	Error    string
}

// GetErrors returns a list of error models by parent id (uuid)
func (a *Access) GetErrors(id uuid.UUID) ([]ErrorModel, error) {
	db := a.getDb()
	result := make([]ErrorModel, 0)
	search := db.Model(&ErrorModel{}).Order("idx")
	if len(id) > 0 {
		search = search.Where("parent_id = ?", id)
	}
	err := search.Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load errors: %w", err)
	}
	return result, err
}

// GetFile returns a file
func (a *Access) GetFile(id uuid.UUID) (*msg.FileMsgItem, error) {
	db := a.getDb()
	var result *msg.FileMsgItem
	search := db.Model(&msg.FileMsgItem{}).Where("id = ?", id)
	err := search.First(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load errors: %w", err)
	}
	return result, err
}

// GetFiles returns a list of files by parent id (uuid)
func (a *Access) GetFiles(id uuid.UUID) ([]msg.FileMsgItem, error) {
	db := a.getDb()
	result := make([]msg.FileMsgItem, 0)
	search := db.Model(&msg.FileMsgItem{}).Order("name")
	if len(id) > 0 {
		search = search.Where("parent_id = ?", id)
	}
	err := search.Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load errors: %w", err)
	}
	return result, err
}

// GetCounters returns a map of counters by parent id (uuid)
func (a *Access) GetCounters(id uuid.UUID) (map[string]string, error) {
	return a.getMap(&counterModel{}, id)
}

// GetStati returns a map of stati by parent id (uuid)
func (a *Access) GetStati(id uuid.UUID) (map[string]string, error) {
	return a.getMap(&statiModel{}, id)
}

func (a *Access) getMap(model any, id uuid.UUID) (map[string]string, error) {
	db := a.getDb()
	result := make(map[string]string)
	list := make([]statiModel, 0)
	search := db.Model(model).Order("name")
	if len(id) > 0 {
		search = search.Where("parent_id = ?", id)
	}
	err := search.Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("cannot load stati: %w", err)
	}
	for _, s := range list {
		result[s.Name] = s.Value
	}
	return result, err
}

// SzenarioModelFromMsg wraps a szenario msg into a model
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

// SaveErrors saves all errors of a szenarion message
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

// SaveStati saves all stati of a szenarion message
func (a *Access) SaveStati(msg *msg.SzenarioEvtMsg) error {
	db := a.getDb()
	var reterr error
	for k, v := range msg.Stati {
		if err := db.Save(statiModel{
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

// SaveCounters saves all counters of a szenarion message
func (a *Access) SaveCounters(msg *msg.SzenarioEvtMsg) error {
	db := a.getDb()
	var reterr error
	for k, v := range msg.Counters {
		if err := db.Save(counterModel{
			ParentID: msg.ID,
			Name:     k,
			Value:    fmt.Sprintf("%v", v),
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

// SaveFiles saves all files of a szenarion message
func (a *Access) SaveFiles(msg *msg.SzenarioEvtMsg) error {
	db := a.getDb()
	var reterr error
	for _, f := range msg.Files {
		f.ParentID = msg.ID
		if f.ID == uuid.Nil {
			f.ID = uuid.New()
		}
		if err := db.Save(&f).Error; err != nil {
			if reterr == nil {
				reterr = err
			} else {
				err = fmt.Errorf("%v %w", reterr, err)
			}
		}
	}
	return reterr
}

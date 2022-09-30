package msg

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

func (m *SzenarioEvtMsg) copy() *SzenarioEvtMsg {
	return &SzenarioEvtMsg{
		ID:         m.ID,
		IncidentID: m.IncidentID,
		Name:       m.Name,
		Time:       m.Time,
		Errors:     m.Errors,
		Counters:   m.Counters,
		Stati:      m.Stati,
		Files:      m.Files,
		Username:   m.Username,
		Region:     m.Region,
		ProbeOS:    m.ProbeOS,
		ProbeHost:  m.ProbeHost,
	}
}

func createSzenarioEvtMsg() *SzenarioEvtMsg {
	e := &SzenarioEvtMsg{
		ID:       uuid.New(),
		Name:     "",
		Region:   viper.GetString(cfg.CheckRegion),
		Time:     time.Time{},
		mu:       sync.RWMutex{},
		Errors:   make([]string, 0),
		Counters: make(map[string]float64),
		Stati:    make(map[string]string),
		Files:    make([]FileMsgItem, 0),
	}
	hostinfo, err := host.Info()
	if err == nil {
		e.ProbeHost = hostinfo.Hostname
		e.ProbeOS = hostinfo.Platform
	}

	return e
}

// NewSzenarioEvtMsg creates a SzenarioEvtMsg
func NewSzenarioEvtMsg(name string, username string, now time.Time) *SzenarioEvtMsg {
	e := createSzenarioEvtMsg()
	e.Name = name
	e.Username = username
	e.Time = now
	return e
}

// SzenarioEvtMsg contains all typicall fields
type SzenarioEvtMsg struct {
	ID         uuid.UUID `json:"ID"`
	IncidentID string    `json:"Incident"`
	Name       string    `json:"Name"`
	Time       time.Time `json:"Time"`
	Username   string    `json:"Username"`
	Region     string    `json:"Region"`
	Errors     []string
	mu         sync.RWMutex
	Counters   map[string]float64 `json:"Counters"`
	Stati      map[string]string
	Files      []FileMsgItem `json:"Files"`
	ProbeOS    string        `json:"OS"`
	ProbeHost  string        `json:"Host"`
}

// AddErr adds a non nil error to the message
func (m *SzenarioEvtMsg) AddErr(e error) {
	if e == nil {
		return
	}
	m.Errors = append(m.Errors, e.Error())
}

// Err returns nil if there are no error and the last non nil error otherwise
func (m *SzenarioEvtMsg) Err() error {
	var err error
	for _, e := range m.Errors {
		if len(e) > 0 {
			err = errors.New(e)
		}
	}
	return err
}

// Errs returns a slice of errors
func (m *SzenarioEvtMsg) Errs() []string {
	return m.Errors
}

// SetStatus sets a status
func (m *SzenarioEvtMsg) SetStatus(key string, val string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Stati[key] = val
}

// SetCounter adds or replaces a counter
// a counter must be int or float
func (m *SzenarioEvtMsg) SetCounter(key string, val float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Counters[key] = val
}

// AddFile adds a file
func (m *SzenarioEvtMsg) AddFile(f *FileMsgItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Files = append(m.Files, *f)
}

// SzenarioEvtMsgFromJSON creates a SzenarioEvtMsg
func SzenarioEvtMsgFromJSON(b []byte) (*SzenarioEvtMsg, error) {
	e := createSzenarioEvtMsg()
	e.AddErr(errors.New(""))
	e.AddFile(&FileMsgItem{})
	err := json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

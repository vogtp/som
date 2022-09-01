package status

import (
	"encoding/json"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
)

// EvtGroup is a status group containing events
type EvtGroup interface {
	Grouper
	AddEvent(evt *msg.SzenarioEvtMsg)
}

// newEvtGroup creates a EvtGroup with a given name and length
// the the length is less than 1 the group is infinite
func newEvtGroup(key string) EvtGroup {
	return &evtGroup{
		Group:  &Group{key: key},
		maxLen: -1,
	}
}

type evtGroup struct {
	*Group
	maxLen int
}

type jsonEvtGrup struct {
	Key      string
	MaxLen   int
	Children []*msg.SzenarioEvtMsg
}

type evtWrapper struct {
	evt *msg.SzenarioEvtMsg
}

// New create a evtWrapper
func (w evtWrapper) New(k string) Grouper {
	return &evtGroup{Group: &Group{key: k}}
}

// Key returns the name of the event
func (w evtWrapper) Key() string {
	return w.evt.ID.String()
}

func (w *evtWrapper) SetConfig(cfg *Config) {
	// to nothing
}

// Level returns the error level
// Down if an error exists OK otherwise
func (w evtWrapper) Level() Level {
	if w.evt.Err() != nil {
		return Down
	}
	return OK
}

// Availability is of the event
// 0 if no error 1 otherwise
func (w evtWrapper) Availability() Availability {
	ret := Availability(1)
	if w.evt.Err() != nil {
		ret = Availability(0)
	}
	return ret
}

// Add panics (not supported for events)
// use AddEvent
func (evtWrapper) Add(Grouper) {
	panic("operation not allowed")
}

// New creates a grouper
func (e *evtGroup) New(k string) Grouper {
	return newEvtGroup(k)
}

// AddEvent adds a event
func (e *evtGroup) AddEvent(evt *msg.SzenarioEvtMsg) {
	e.Add(&evtWrapper{evt})
}

// Key returns the name or key of the group
func (e *evtGroup) Key() string {
	return e.key
}

// Add a sub grouper
func (e *evtGroup) Add(c Grouper) {
	e.Group.Add(c)
	if e.maxLen < 0 {
		e.maxLen = viper.GetInt(cfg.AlertIncidentCorrelationEvents)
	}
	if e.maxLen > 0 && len(e.children) > e.maxLen {
		e.children = e.children[1:]
	}
}

// Level returns the error level
func (e evtGroup) Level() Level {
	lcIdx := len(e.children) - 1
	if lcIdx < 0 {
		return Unknown
	}
	lastChild := e.Group.children[lcIdx]
	if e.cfg.UnknownTimeout > 0 {
		if evt, ok := lastChild.(*evtWrapper); ok {
			// if the newest event is older than 100 * repeat time set the status to unknown
			if time.Since(evt.evt.Time) > e.cfg.UnknownTimeout {
				return Unknown
			}
		}
	}
	lvl := e.Group.Level()
	if lvl == OK {
		return lvl
	}
	if lastChild.Level() == OK {
		lvl--
	}
	return lvl
}

// MarshalJSON ...
func (e evtGroup) MarshalJSON() ([]byte, error) {
	j := &jsonEvtGrup{
		Key:      e.key,
		MaxLen:   e.maxLen,
		Children: make([]*msg.SzenarioEvtMsg, 0),
	}
	for _, c := range e.children {
		if evt, ok := c.(*evtWrapper); ok {
			j.Children = append(j.Children, evt.evt)
		}
	}
	return json.Marshal(j)
}

// UnmarshalJSON ...
func (e evtGroup) UnmarshalJSON(data []byte) error {
	j := &jsonEvtGrup{
		Children: make([]*msg.SzenarioEvtMsg, 0),
	}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	e.key = j.Key
	e.maxLen = j.MaxLen
	for _, evt := range j.Children {
		e.Add(&evtWrapper{evt: evt})
	}
	return nil
}

// MarshalJSON ...
func (w evtWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.evt)
}

// UnmarshalJSON ...
func (w evtWrapper) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, w.evt)
}

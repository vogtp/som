//go:generate stringer -type=IncidentMsgType

package msg

import "time"

// IncidentMsgType idicates what type of msg it is
type IncidentMsgType int

const (
	// OpenIncident a new incident
	OpenIncident IncidentMsgType = iota
	// UpdateIncident a existing incident
	UpdateIncident
	// CloseIncident the incident
	CloseIncident
)

// IncidentMsg communicates the state of the incident
type IncidentMsg struct {
	*SzenarioEvtMsg
	Type      IncidentMsgType
	Start     time.Time
	End       time.Time
	IntLevel  int    `json:"Level"`
	ByteState []byte `json:"State"`
}

// NewIncidentMsg creates a incident message
func NewIncidentMsg(t IncidentMsgType, m *SzenarioEvtMsg) *IncidentMsg {
	msg := &IncidentMsg{
		Type:           t,
		SzenarioEvtMsg: m.copy(),
	}
	switch t {
	case OpenIncident:
		msg.Start = m.Time
	case CloseIncident:
		msg.End = m.Time
	}
	return msg
}

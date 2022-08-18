package msg

// AlertMsg is a alert
type AlertMsg struct {
	*SzenarioEvtMsg
	Level string
}

// NewAlert converts the Event to an Alert
func NewAlert(m *SzenarioEvtMsg) *AlertMsg {
	return &AlertMsg{
		SzenarioEvtMsg: m.copy(),
	}
}

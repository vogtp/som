package bus

import (
	"fmt"
)

func (m *Manager) Emit(subject string, data []byte) error {
	if err := m.EnsureConnected(); err != nil {
		return fmt.Errorf("ensuring connected: %w", err)
	}
	sl := m.slog.With("subject", subject)

	err := m.conn.Publish(subject, data)

	if err != nil {
		return fmt.Errorf("publish messsage to %s: %w", subject, err)
	}

	sl.Debug("Sent message", "msg", string(data))
	return nil
}

func (m *Manager) Ask(subject string, data []byte) (*Message, error) {
	if err := m.EnsureConnected(); err != nil {
		return nil, fmt.Errorf("ensuring connected: %w", err)
	}
	//sl := m.slog.With("subject", subject, "bus", "ask")

	resp, err := m.conn.Request(subject, data, m.timeout)
	if err != nil {
		return nil, fmt.Errorf("requesting answer: %w", err)
	}

	busMsg := &Message{
		Subject: resp.Subject,
		Body:    resp.Data,
	}
	return busMsg, nil
}

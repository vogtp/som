package bus

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/vogtp/som/pkg/core/log"
)

type ReceiveFunc func(subject string, msg *Message)

func (m *Manager) Subscribe(subject string, recFunc ReceiveFunc) (unsubscribeFunc, error) {
	if err := m.EnsureConnected(); err != nil {
		return nil, fmt.Errorf("ensuring connected: %w", err)
	}
	sl := m.slog.With("subject", subject, "bus", "receive")

	sub, err := m.conn.Subscribe(subject, func(msg *nats.Msg) {
		m := Message{
			Subject: msg.Subject,
			Body:    msg.Data,
		}
		sl.Debug("Got messgae", "subject", subject)
		recFunc(msg.Subject, &m)
	})

	return func() { _ = sub.Unsubscribe() }, err
}

type AnswerFunc func(subject string, msg *Message) ([]byte, error)

func (m *Manager) Respond(subject string, answerFunc AnswerFunc) (unsubscribeFunc, error) {
	if err := m.EnsureConnected(); err != nil {
		return nil, fmt.Errorf("ensuring connected: %w", err)
	}
	sl := m.slog.With("subject", subject, "bus", "receive")

	sub, err := m.conn.Subscribe(subject, func(msg *nats.Msg) {
		busMsg := Message{
			Subject: msg.Subject,
			Body:    msg.Data,
		}
		sl.Debug("Got messgae", "subject", subject)
		d, err := answerFunc(msg.Subject, &busMsg)
		if err != nil {
			sl.Warn("Bus respond got error", log.Error, err, "data", string(d))
		}

		if err := msg.Respond(d); err != nil {
			sl.Warn("Failure responding to message", log.Error, err)
		}
	})

	return func() { _ = sub.Unsubscribe() }, err
}

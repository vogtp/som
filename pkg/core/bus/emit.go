package bus

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func (m *manager) Emit(subject string, data []byte) error {
	sl := m.slog.With("subject", subject)

	err := m.conn.Publish(subject, data)

	if err != nil {
		return fmt.Errorf("publish messsage to %s: %w", subject, err)
	}

	sl.Debug("Sent message", "msg", string(data))
	return nil
}

func (m *manager) Ask(subject string, data []byte) (*Message, error) {
	sl := m.slog.With("subject", subject, "bus", "ask")
	replySubject := fmt.Sprintf("%s.reply", subject)
	msg := nats.Msg{
		Subject: subject,
		Data:    data,
		Reply:   replySubject,
	}
	sl.Debug("Sending message", "replySubject", replySubject, "data", string(data))
	if err := m.conn.PublishMsg(&msg); err != nil {
		return nil, fmt.Errorf("sending ask message: %w", err)
	}

	sl.Debug("Subscribe to replies", "replySubject", replySubject)
	sub, err := m.conn.SubscribeSync(replySubject)
	if err != nil {
		return nil, fmt.Errorf("subscribing to ask replies: %w", err)
	}
	defer func() { _ = sub.Unsubscribe() }()
	ansMsg, err := sub.NextMsg(m.timeout)
	if err != nil {
		return nil, fmt.Errorf("getting reply message: %w", err)
	}
	busMsg := &Message{
		Subject: subject,
		Body:    ansMsg.Data,
	}
	return busMsg, nil
}

package bus

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

func (m *manager) Emit(ctx context.Context, routingKey string, data []byte) error {
	sl := m.slog.With("routingKey", routingKey)

	err := m.conn.Publish(routingKey, data)

	if err != nil {
		return fmt.Errorf("publish messsage to %s: %w", routingKey, err)
	}

	sl.Debug("Sent message", "msg", string(data))
	return nil
}

func (m *manager) Ask(ctx context.Context, routingKey string, data []byte) (*Message, error) {
	sl := m.slog.With("routingKey", routingKey, "bus", "ask")
	replySubject := fmt.Sprintf("%s.reply", routingKey)
	msg := nats.Msg{
		Subject: routingKey,
		Data:    data,
		Reply:   replySubject,
	}
	sl.Info("Sending message", "replySubject", replySubject, "data", string(data))
	if err := m.conn.PublishMsg(&msg); err != nil {
		return nil, fmt.Errorf("sending ask message: %w", err)
	}

	m.conn.Subscribe(replySubject, func(msg *nats.Msg) {
		sl.Info("reply subscibe ****", "data", msg.Data)
	})

	sl.Info("Sub to replies", "replySubject", replySubject)

	sub, err := m.conn.SubscribeSync(replySubject)
	if err != nil {
		return nil, fmt.Errorf("subscribing to ask replies: %w", err)
	}
	defer sub.Unsubscribe()
	ansMsg, err := sub.NextMsg(m.timeout)
	if err != nil {
		return nil, fmt.Errorf("getting reply message: %w", err)
	}
	busMsg := &Message{
		RoutingKey: routingKey,
		Body:       ansMsg.Data,
	}
	return busMsg, nil
}

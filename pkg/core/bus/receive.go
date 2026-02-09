package bus

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/vogtp/som/pkg/core/log"
)

type ReceiveFunc func(routingKey string, msg *Message)

func (m *manager) Receive(ctx context.Context, routingKey string, recFunc ReceiveFunc) (unsubscribecloseFunc, error) {
	sl := m.slog.With("routingKey", routingKey, "bus", "receive")
	// ctx, cancel := context.WithTimeout(ctx, m.timeout)
	// defer cancel()

	sub, err := m.conn.Subscribe(routingKey, func(msg *nats.Msg) {
		m := Message{
			RoutingKey: msg.Subject,
			Body:       msg.Data,
		}
		sl.Debug("Got messgae", "routingKey", routingKey)
		recFunc(msg.Subject, &m)
	})

	return func() { sub.Unsubscribe() }, err
}

type AnswerFunc func(routingKey string, msg *Message) ([]byte, error)

func (m *manager) Answer(ctx context.Context, routingKey string, answerFunc AnswerFunc) (unsubscribecloseFunc, error) {
	sl := m.slog.With("routingKey", routingKey, "bus", "receive")
	// ctx, cancel := context.WithTimeout(ctx, m.timeout)
	// defer cancel()

	sub, err := m.conn.Subscribe(routingKey, func(msg *nats.Msg) {
		busMsg := Message{
			RoutingKey: msg.Subject,
			Body:       msg.Data,
		}
		sl.Debug("Got messgae", "routingKey", routingKey)
		d, err := answerFunc(msg.Subject, &busMsg)
		if err != nil {
			sl.Warn("Bus answer got error", log.Error, err, "data", string(d))
		}
		reply := nats.Msg{
			Subject: msg.Reply,
			Data:    d,
		}
		if err := m.conn.PublishMsg(&reply); err != nil {
			sl.Warn("Failure publishing reply message", log.Error, err, "msg", reply)
		}
	})

	return func() { sub.Unsubscribe() }, err
}

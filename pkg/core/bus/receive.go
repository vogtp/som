package bus

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/vogtp/som/pkg/core/log"
)

type ReceiveFunc func(subject string, msg *Message)

func (m *manager) Receive(ctx context.Context, subject string, recFunc ReceiveFunc) (unsubscribecloseFunc, error) {
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

func (m *manager) Answer(ctx context.Context, subject string, answerFunc AnswerFunc) (unsubscribecloseFunc, error) {
	sl := m.slog.With("subject", subject, "bus", "receive")


	sub, err := m.conn.Subscribe(subject, func(msg *nats.Msg) {
		busMsg := Message{
			Subject: msg.Subject,
			Body:    msg.Data,
		}
		sl.Debug("Got messgae", "subject", subject)
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

	return func() { _ = sub.Unsubscribe() }, err
}

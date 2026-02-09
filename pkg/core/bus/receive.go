package bus

import (
	"context"

	"github.com/nats-io/nats.go"
)

type ReceiveFunc func(routingKey string, msg *Message)

func (m *manager) Receive(ctx context.Context, routingKey string, recFunc ReceiveFunc) error {
	sl := m.slog.With("routingKey", routingKey, "bus", "receive")
	// ctx, cancel := context.WithTimeout(ctx, m.timeout)
	// defer cancel()

	_, err := m.conn.Subscribe(routingKey, func(msg *nats.Msg) {
		m := Message{
			RoutingKey: msg.Subject,
			Body:       msg.Data,
		}
		sl.Debug("Got messgae", "routingKey", routingKey)
		recFunc(msg.Subject, &m)
	})

	return err
}

type AnswerFunc func(routingKey string, msg Message) ([]byte, error)

func (m *manager) Answer(ctx context.Context, routingKey string, answerFunc AnswerFunc) error {
	// sl := m.slog.With("routingKey", routingKey, "bus", "answer")
	// q, err := m.channel.QueueDeclare(
	// 	"rpc_queue", // name
	// 	false,       // durable
	// 	false,       // delete when unused
	// 	false,       // exclusive
	// 	false,       // no-wait
	// 	nil,         // arguments
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to declare a queue: %w", err)
	// }

	// err = m.channel.Qos(
	// 	1,     // prefetch count
	// 	0,     // prefetch size
	// 	false, // global
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to set QoS: %w", err)
	// }
	// msgs, err := m.channel.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	false,  // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,    // args
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to register a consumer: %w", err)
	// }

	// go func() {
	// 	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	// 	defer cancel()
	// 	for {
	// 		select {
	// 		case d := <-msgs:
	// 			r, err := answerFunc(routingKey, d)
	// 			if err != nil {
	// 				sl.Warn("Answer func yielded an error", "err", err)
	// 			}

	// 			err = m.channel.PublishWithContext(ctx,
	// 				"",        // exchange
	// 				d.ReplyTo, // routing key
	// 				false,     // mandatory
	// 				false,     // immediate
	// 				amqp.Publishing{
	// 					ContentType:   "text/plain",
	// 					CorrelationId: d.CorrelationId,
	// 					Body:          r,
	// 				})
	// 			if err != nil {
	// 				sl.Warn("Publish response", "err", err)
	// 				continue
	// 			}

	// 			if err := d.Ack(false); err != nil {
	// 				sl.Warn("Sending ack", "err", err)
	// 			}
	// 		case <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()
	// return err
	return nil
}

// func (m *manager) Answer(ctx context.Context, routingKey string, answerFunc AnswerFunc) error {
// 	sl := m.slog.With("routingKey", routingKey, "bus", "answer")
// 	q, err := m.channel.QueueDeclare(
// 		"rpc_queue", // name
// 		false,       // durable
// 		false,       // delete when unused
// 		false,       // exclusive
// 		false,       // no-wait
// 		nil,         // arguments
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to declare a queue: %w", err)
// 	}

// 	err = m.channel.Qos(
// 		1,     // prefetch count
// 		0,     // prefetch size
// 		false, // global
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to set QoS: %w", err)
// 	}

// 	msgs, err := m.channel.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		false,  // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to register a consumer: %w", err)
// 	}

// 	go func() {
// 		ctx, cancel := context.WithTimeout(ctx, m.timeout)
// 		defer cancel()
// 		for {
// 			var d amqp.Delivery
// 			select {
// 			case d = <-msgs:
// 				sl.Debug("Got message", "msg.routingKey", d.MessageId, "correlationId", d.CorrelationId, "payload", string(d.Body))
// 				if !strings.EqualFold(routingKey, d.MessageId) {
// 					continue
// 				}
// 			case <-ctx.Done():
// 				err = fmt.Errorf("timeout %s reached: %w", m.timeout, ctx.Err())
// 				return
// 			}
// 			resp, err := answerFunc(routingKey, d)
// 			if err != nil {
// 				sl.Warn("Bus answer func returned an error", log.Error, err)
// 				continue
// 			}

// 			err = m.channel.PublishWithContext(ctx,
// 				"",        // exchange
// 				d.ReplyTo, // routing key
// 				false,     // mandatory
// 				false,     // immediate
// 				amqp.Publishing{
// 					ContentType:   "text/plain",
// 					CorrelationId: d.CorrelationId,
// 					MessageId:     routingKey,
// 					Body:          resp,
// 				})
// 			if err != nil {
// 				sl.Warn("Sending bus answer returned an error", log.Error, err)
// 				continue
// 			}

// 			if err := d.Ack(false); err != nil {
// 				sl.Warn("Bus answer ack got error", log.Error, err)
// 			}
// 		}
// 	}()
// 	return err
// }

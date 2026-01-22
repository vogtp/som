package bus

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vogtp/som/pkg/core/log"
)

type ReceiveFunc func(amqp.Delivery)

func (m *manager) Receive(routingKey string, recFunc ReceiveFunc) error {
	sl := m.slog.With("routingKey", routingKey)
	q, err := m.channel.QueueDeclare(
		fmt.Sprintf("%s.queue", routingKey), // name
		true,                                // durable
		false,                               // delete when unused
		true,                                // exclusive
		false,                               // no-wait
		nil,                                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = m.channel.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		somTopic,   // exchange
		false,
		nil)
	if err != nil {
		return fmt.Errorf("binding queue to exchange %s: %w", somTopic, err)
	}
	sl.Info("Binding queue to bus")

	msgs, err := m.channel.Consume(
		q.Name,                                // queue
		fmt.Sprintf("%s.consume", routingKey), // consumer
		true,                                  // auto ack
		false,                                 // exclusive
		false,                                 // no local
		false,                                 // no wait
		nil,                                   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer for %s: %w", routingKey, err)
	}

	go func() {
		for m := range msgs {
			sl.Debug("Received message", "msg", m)
			recFunc(m)
		}
	}()
	return nil
}

type AnswerFunc func(amqp.Delivery) ([]byte, error)

func (m *manager) Answer(ctx context.Context, routingKey string, answerFunc AnswerFunc) error {
	routingKey = fmt.Sprintf("%s.askanswer", routingKey)
	sl := m.slog.With("routingKey", routingKey)
	q, err := m.channel.QueueDeclare(
		routingKey, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = m.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := m.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, m.timeout)
		defer cancel()
		for d := range msgs {
			//FIXME handle context
			resp, err := answerFunc(d)
			if err != nil {
				sl.Warn("Bus answer func returned an error", log.Error, err)
				continue
			}

			err = m.channel.PublishWithContext(ctx,
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          resp,
				})
			if err != nil {
				sl.Warn("Sending bus answer returned an error", log.Error, err)
				continue
			}

			if err := d.Ack(false); err != nil {
				sl.Warn("Bus answer ack got error", log.Error, err)
			}
		}
	}()
	return nil
}

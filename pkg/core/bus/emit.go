package bus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (m *manager) Emit(ctx context.Context, routingKey string, data []byte) error {
	sl := m.slog.With("routingKey", routingKey)
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	err := m.channel.PublishWithContext(ctx,
		somTopic,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			// DeliveryMode: 2, // persisten
			Body:      data,
			MessageId: uuid.NewString(),
			// ReplyTo: ,
		})
	if err != nil {
		return fmt.Errorf("publish messsage to %s: %w", routingKey, err)
	}

	sl.Debug("Sent message", "msg", string(data))
	return nil
}

func (m *manager) Ask(ctx context.Context, routingKey string, data []byte) (*amqp.Delivery, error) {
	routingKey = fmt.Sprintf("%s.askanswer", routingKey)
	q, err := m.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("generating reply queue %s: %w", fmt.Sprintf("%s.reply", routingKey), err)
	}

	msgs, err := m.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("register a consumer for reply queue %s: %w", fmt.Sprintf("%s.reply", routingKey), err)
	}

	corrId := uuid.NewString()

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	err = m.channel.PublishWithContext(ctx,
		"",   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          data,
		})
	if err != nil {
		return nil, fmt.Errorf("publish messsage to %s: %w", routingKey, err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout %s reached", m.timeout)
		case d := <-msgs:
			if corrId == d.CorrelationId {
				return &d, nil
			}
		}
	}

}

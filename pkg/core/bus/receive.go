package bus

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ReceiveFunc func(amqp.Delivery)

func (m *manager) Receive(routingKey string, recFunc ReceiveFunc) error {
	sl := m.slog.With("routingKey", routingKey)
	q, err := m.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
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
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
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

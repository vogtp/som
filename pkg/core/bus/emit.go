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
			ContentType:  "text/plain",
			// DeliveryMode: 2, // persisten
			Body:         data,
			MessageId: uuid.NewString(),
		})
	if err != nil {
		return fmt.Errorf("publish messsage to %s: %w", routingKey, err)
	}

	sl.Debug("Sent message", "msg", string(data))
	return nil
}

package bus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	somTopic = "som_topic"
)

type Manager interface {
	Emit(ctx context.Context, routingKey string, data []byte) error
	Receive(routingKey string, recFunc ReceiveFunc) error
	Close() //Close all AMQP resouces
}

type manager struct {
	slog *slog.Logger

	timeout time.Duration

	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(slog *slog.Logger) (Manager, error) {
	conn, err := amqp.Dial(cfg.AmqpURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the AMQP server %q: %w", cfg.AmqpURL(), err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("creating the AMQP channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		somTopic, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("creating the %q topic: %w", somTopic, err)
	}
	m := &manager{
		slog:    slog.With(log.Component, "bus"),
		timeout: 5*time.Second,
		conn:    conn,
		channel: ch,
	}
	return m, nil
}

func (m *manager) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
	if m.channel != nil {
		m.channel.Close()
	}
}

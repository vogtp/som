package bus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/vogtp/som/pkg/core/log"
)

const (
	somTopic = "som.topic"
)

type Manager interface {
	Emit(ctx context.Context, routingKey string, data []byte) error
	Receive(ctx context.Context, routingKey string, recFunc ReceiveFunc) error

	Ask(ctx context.Context, routingKey string, data []byte) (*amqp.Delivery, error)
	Answer(ctx context.Context, routingKey string, answerFunc AnswerFunc) error

	Close() //Close all AMQP resouces
}

type manager struct {
	slog *slog.Logger

	timeout time.Duration
	client  *redis.Client
}

func New(ctx context.Context, slog *slog.Logger) (Manager, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := client.ClientInfo(ctx).Err(); err != nil {
		return nil, fmt.Errorf("create redis client: %w", err)
	}

	m := manager{
		slog:    slog.With("bus", "redis"),
		client:  client,
		timeout: 5 * time.Second,
	}
	return &m, nil
}

func (m *manager) SetTimeout(d time.Duration) {
	if d < time.Millisecond {
		m.slog.Warn("Not setting bus timeout since it is too low", "timeout", d, log.Stacktrace())
	}
	m.timeout = d
}

func (m *manager) Close() {
	if m.client != nil {
		m.client.Close()
	}
}

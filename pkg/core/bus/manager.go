package bus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vogtp/som/pkg/core/log"
)

const (
	somTopic = "som.topic"
	natsURL  = "nats://0.0.0.0:4222" //TODO move to config
)

type Manager interface {
	Emit(ctx context.Context, subject string, data []byte) error
	Receive(ctx context.Context, subject string, recFunc ReceiveFunc) (unsubscribecloseFunc, error)

	Ask(ctx context.Context, subject string, data []byte) (*Message, error)
	Answer(ctx context.Context, subject string, answerFunc AnswerFunc) (unsubscribecloseFunc, error)

	Close() //Close all resouces
}

type unsubscribecloseFunc func()

type manager struct {
	slog *slog.Logger

	timeout time.Duration
	conn    *nats.Conn
}

func New(ctx context.Context, slog *slog.Logger) (Manager, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("connect to nats server: %w", err)
	}

	m := manager{
		slog:    slog.With("bus", "redis"),
		conn:    conn,
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
	if m.conn != nil {
		m.conn.Close()
	}
}

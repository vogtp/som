package bus

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/msgtype"
)

type unsubscribeFunc func()

type Manager struct {
	slog *slog.Logger

	timeout time.Duration
	conn    *nats.Conn

	Szenario *eventHandler[msg.SzenarioEvtMsg]
	Alert    *eventHandler[msg.AlertMsg]
	Incident *eventHandler[msg.IncidentMsg]
}

func New(slog *slog.Logger) (*Manager, error) {

	m := Manager{
		slog:    slog.With(log.Component, "bus"),
		timeout: viper.GetDuration(cfg.BusTimeout),
	}
	if err := m.EnsureConnected(); err != nil {
		return nil, fmt.Errorf("connect to nats server: %w", err)
	}
	m.initEventHandlers()
	return &m, nil
}

func (m *Manager) EnsureConnected() error {
	url := viper.GetString(cfg.BusURL)
	m.slog = m.slog.With("nats.url", url)
	if m.conn == nil {
		conn, err := nats.Connect(url)
		if err != nil {
			return err
		}
		m.conn = conn
	}
	if err := m.conn.Flush(); err == nil {
		return nil
	}
	conn, err := nats.Connect(url)
	if err != nil {
		return err
	}
	m.conn = conn
	return nil
}

// newBus creates a new eventbus
func (m *Manager) initEventHandlers() {
	m.Szenario = newHandler[msg.SzenarioEvtMsg](
		m.slog,
		m,
		msgtype.Event,
	)
	m.Alert = newHandler[msg.AlertMsg](
		m.slog,
		m,
		msgtype.Alert,
	)
	m.Incident = newHandler[msg.IncidentMsg](
		m.slog,
		m,
		msgtype.Incident,
	)
}

func (m *Manager) SetTimeout(d time.Duration) {
	if d < time.Millisecond {
		m.slog.Warn("Not setting bus timeout since it is too low", "timeout", d, log.Stacktrace())
	}
	m.timeout = d
}

func (m *Manager) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
	if m.Incident != nil {
		// if Incident exists all exist
		m.Szenario.cleanup()
		m.Alert.cleanup()
		m.Incident.cleanup()
	}
}

// WaitMsgProcessed waits until the managed cannels have their messages sent
func (m *Manager) WaitMsgProcessed() {
	m.Szenario.WaitMsgProcessed()
	m.Alert.WaitMsgProcessed()
	m.Incident.WaitMsgProcessed()
}

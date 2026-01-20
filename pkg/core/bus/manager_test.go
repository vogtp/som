package bus_test

import (
	"log/slog"
	"strings"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/bus"
	"github.com/vogtp/som/pkg/core/cfg"
)

const (
	somAmqpPass = "somMsgBusPassword"
)

func TestNewBusManager(t *testing.T) {
	cfg.Parse()
	_, err := bus.New(slog.Default())
	if err == nil {
		t.Errorf("AMQP bus works without password")
	}
	viper.Set(cfg.AmqpPasswort, somAmqpPass)
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Errorf("Initalise AMQP bus: %v", err)
	}
	if m == nil {
		t.Errorf("AMQP bus is nil")
	}
	m.Close()
}

func TestSendReceive(t *testing.T) {
	cfg.Parse()
	viper.Set(cfg.AmqpPasswort, somAmqpPass)
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Fatalf("Initalise AMQP bus: %v", err)
	}
	recMsg := ""
	var recTime time.Time
	rk := "som.testing"
	wait := make(chan any)
	err = m.Receive(rk, func(d amqp.Delivery) {
		recMsg = string(d.Body)
		recTime = time.Now()
		close(wait)
	})
	if err != nil {
		t.Errorf("cannot receive: %v", err)
	}
	msg := "Test message"
	err = m.Emit(t.Context(), rk, []byte(msg))
	sendTime := time.Now()
	if err != nil {
		t.Errorf("cannot emit: %v", err)
	}
	<-wait
	if !strings.EqualFold(msg, recMsg) {
		t.Fatalf("Got message %s expected %s", recMsg, msg)
	}
	if recTime.Sub(sendTime) > time.Millisecond {
		t.Errorf("sending took too long: %s", recTime.Sub(sendTime))
	}
}

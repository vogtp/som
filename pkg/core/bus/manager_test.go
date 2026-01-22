package bus_test

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/bus"
	"github.com/vogtp/som/pkg/core/cfg"
)

func TestNewBusManager(t *testing.T) {
	cfg.Parse()
	pw := viper.Get(cfg.AmqpPasswort)
	viper.Set(cfg.AmqpPasswort, "")
	_, err := bus.New(slog.Default())
	if err == nil {
		t.Errorf("AMQP bus works without password")
	}
	viper.Set(cfg.AmqpPasswort, pw)
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
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Fatalf("Initalise AMQP bus: %v", err)
	}
	recMsg := ""
	var recTime time.Time
	rk := "som.testing"
	wait := make(chan any)
	err = m.Receive(rk, func(r string, d amqp.Delivery) {
		if !strings.EqualFold(r, rk) {
			t.Errorf("Routing keys do not match: have: %s want: %s", r, rk)
		}
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

func TestAskAnswer(t *testing.T) {
	cfg.Parse()
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Fatalf("Initalise AMQP bus: %v", err)
	}
	ansRecMsg := ""
	ansSendMsg := ""
	var recTime time.Time
	rk := "som.testing"
	wait := make(chan any)
	err = m.Answer(t.Context(), rk, func(r string, d amqp.Delivery) ([]byte, error) {
		if !strings.EqualFold(r, rk) {
			t.Errorf("Routing keys do not match: have: %s want: %s", r, rk)
		}
		ansRecMsg = string(d.Body)
		ansSendMsg = fmt.Sprintf("answer-%s", ansRecMsg)
		recTime = time.Now()
		close(wait)
		return []byte(ansSendMsg), nil
	})
	if err != nil {
		t.Errorf("cannot answer: %v", err)
	}
	msg := "Test message"
	resp, err := m.Ask(t.Context(), rk, []byte(msg))
	sendTime := time.Now()
	if err != nil {
		t.Fatalf("cannot ask: %v", err)
	}
	<-wait
	if !strings.EqualFold(msg, ansRecMsg) {
		t.Fatalf("Got message %s expected %s", ansRecMsg, msg)
	}
	if !strings.EqualFold(ansSendMsg, string(resp.Body)) {
		t.Fatalf("Got answer %s expected %s", ansSendMsg, resp.Body)
	}
	if recTime.Sub(sendTime) > time.Millisecond {
		t.Errorf("sending took too long: %s", recTime.Sub(sendTime))
	}
}

package bus_test

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/vogtp/som/pkg/core/bus"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
)

type timeouter interface {
	SetTimeout(time.Duration)
}

func TestNewBusManager(t *testing.T) {
	cfg.Parse()
	// pw := viper.Get(cfg.AmqpPasswort)
	// viper.Set(cfg.AmqpPasswort, "")
	// _, err := bus.New(slog.Default())
	// if err == nil {
	// 	t.Errorf("AMQP bus works without password")
	// }
	// viper.Set(cfg.AmqpPasswort, pw)
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Errorf("Initalise bus: %v", err)
	}
	if m == nil {
		t.Errorf("AMQP bus is nil")
	}
	m.Close()
}

func TestSendReceive(t *testing.T) {
	cfg.Parse()
	log.Level.Set(slog.LevelDebug)
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Fatalf("Initalise AMQP bus: %v", err)
	}
	m.SetTimeout(time.Minute * 5)
	recMsg := ""
	var recTime time.Time
	rk := "som.testing"
	wait := make(chan any, 1)
	defer close(wait)
	close, err := m.Receive(rk, func(r string, m *bus.Message) {
		if !strings.EqualFold(r, rk) {
			t.Errorf("Routing keys do not match: have: %s want: %s", r, rk)
		}
		recMsg = string(m.Body)
		recTime = time.Now()
		wait <- 1
	})
	if err != nil {
		t.Errorf("cannot receive: %v", err)
	}
	defer close()
	msg := "Test message Send/Receive"
	err = m.Emit(rk, []byte(msg))
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
	log.Level.Set(slog.LevelDebug)
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Fatalf("Initalise AMQP bus: %v", err)
	}
	ansRecMsg := ""
	ansSendMsg := ""
	var recTime time.Time
	rk := "som.testing"
	close, err := m.Answer(rk, func(r string, d *bus.Message) ([]byte, error) {
		if !strings.EqualFold(r, rk) {
			t.Errorf("Routing keys do not match: have: %s want: %s", r, rk)
		}
		ansRecMsg = string(d.Body)
		ansSendMsg = fmt.Sprintf("answer-%s", ansRecMsg)
		recTime = time.Now()
		return []byte(ansSendMsg), nil
	})
	if err != nil {
		t.Errorf("cannot answer: %v", err)
	}
	defer close()
	msg := "Test message Ask/Answer"
	resp, err := m.Ask(rk, []byte(msg))
	sendTime := time.Now()
	if err != nil {
		t.Fatalf("cannot ask: %v", err)
	}
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

func TestAskAnswerWildcard(t *testing.T) {
	tests := []struct {
		subjectSend string
		subjectRec  string
	}{
		{"som.testing.test", "som.testing.test"},
		{"som.testing.test", "som.testing.*"},
	}

	cfg.Parse()
	log.Level.Set(slog.LevelDebug)
	m, err := bus.New(slog.Default())
	if err != nil {
		t.Fatalf("Initalise AMQP bus: %v", err)
	}
	m.SetTimeout(time.Second * 5)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.subjectSend, tt.subjectRec), func(t *testing.T) {

			ansRecMsg := ""
			ansSendMsg := ""
			var recTime time.Time
			close, err := m.Answer(tt.subjectRec, func(r string, d *bus.Message) ([]byte, error) {
				if !strings.EqualFold(r, tt.subjectSend) {
					t.Errorf("Routing keys do not match: have: %s want: %s", r, tt.subjectSend)
				}
				ansRecMsg = string(d.Body)
				ansSendMsg = fmt.Sprintf("answer-%s", ansRecMsg)
				recTime = time.Now()
				return []byte(ansSendMsg), nil
			})
			if err != nil {
				t.Errorf("cannot answer: %v", err)
			}
			defer close()
			msg := "Test message Ask/Answer WildCard"
			sendTime := time.Now()
			resp, err := m.Ask(tt.subjectSend, []byte(msg))
			if err != nil {
				t.Fatalf("cannot ask: %v", err)
			}
			if !strings.EqualFold(msg, ansRecMsg) {
				t.Fatalf("Got message %s expected %s", ansRecMsg, msg)
			}
			if !strings.EqualFold(ansSendMsg, string(resp.Body)) {
				t.Fatalf("Got answer %s expected %s", ansSendMsg, resp.Body)
			}
			if recTime.Sub(sendTime) > time.Millisecond {
				t.Errorf("sending took too long: %s", recTime.Sub(sendTime))
			}
		})
	}
}

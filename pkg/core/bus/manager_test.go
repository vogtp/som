package bus_test

import (
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/bus"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
)

type timeouter interface {
	SetTimeout(time.Duration)
}

func TestNewBusManager(t *testing.T) {
	cfg.Parse()
	pw := viper.Get(cfg.AmqpPasswort)
	// viper.Set(cfg.AmqpPasswort, "")
	// _, err := bus.New(t.Context(),slog.Default())
	// if err == nil {
	// 	t.Errorf("AMQP bus works without password")
	// }
	viper.Set(cfg.AmqpPasswort, pw)
	m, err := bus.New(t.Context(), slog.Default())
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
	log.Level.Set(slog.LevelDebug)
	m, err := bus.New(t.Context(), slog.Default())
	if err != nil {
		t.Fatalf("Initalise AMQP bus: %v", err)
	}
	if to, ok := m.(timeouter); ok {
		to.SetTimeout(time.Minute * 5)
	}
	recMsg := ""
	var recTime time.Time
	rk := "som.testing"
	wait := make(chan any)
	err = m.Receive(t.Context(), rk, func(r string, m bus.Message) {
		if !strings.EqualFold(r, rk) {
			t.Errorf("Routing keys do not match: have: %s want: %s", r, rk)
		}
		recMsg = string(m.Payload)
		recTime = time.Now()
		close(wait)
	})
	if err != nil {
		t.Errorf("cannot receive: %v", err)
	}
	msg := "Test message Send/Receive"
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

// func TestAskAnswer(t *testing.T) {
// 	cfg.Parse()
// 	log.Level.Set(slog.LevelDebug)
// 	m, err := bus.New(slog.Default())
// 	if err != nil {
// 		t.Fatalf("Initalise AMQP bus: %v", err)
// 	}
// 	ansRecMsg := ""
// 	ansSendMsg := ""
// 	var recTime time.Time
// 	rk := "som.testing"
// 	wait := make(chan int)
// 	err = m.Answer(t.Context(), rk, func(r string, d amqp.Delivery) ([]byte, error) {
// 		if !strings.EqualFold(r, rk) {
// 			t.Errorf("Routing keys do not match: have: %s want: %s", r, rk)
// 		}
// 		ansRecMsg = string(d.Body)
// 		ansSendMsg = fmt.Sprintf("answer-%s", ansRecMsg)
// 		recTime = time.Now()
// 		wait <- 1
// 		return []byte(ansSendMsg), nil
// 	})
// 	if err != nil {
// 		t.Errorf("cannot answer: %v", err)
// 	}
// 	msg := "Test message Ask/Answer"
// 	resp, err := m.Ask(t.Context(), rk, []byte(msg))
// 	sendTime := time.Now()
// 	if err != nil {
// 		t.Fatalf("cannot ask: %v", err)
// 	}
// 	<-wait
// 	if !strings.EqualFold(msg, ansRecMsg) {
// 		t.Fatalf("Got message %s expected %s", ansRecMsg, msg)
// 	}
// 	if !strings.EqualFold(ansSendMsg, string(resp.Body)) {
// 		t.Fatalf("Got answer %s expected %s", ansSendMsg, resp.Body)
// 	}
// 	if recTime.Sub(sendTime) > time.Millisecond {
// 		t.Errorf("sending took too long: %s", recTime.Sub(sendTime))
// 	}
// }

// func TestAskAnswerWildcard(t *testing.T) {
// 	tests := []struct {
// 		routingKeySend string
// 		routingKeyRec  string
// 	}{
// 		{"som.testing.test", "som.testing.test"},
// 		//	{"som.testing.test", "som.testing.test"},
// 	}

// 	cfg.Parse()
// 	log.Level.Set(slog.LevelDebug)
// 	m, err := bus.New(slog.Default())
// 	if err != nil {
// 		t.Fatalf("Initalise AMQP bus: %v", err)
// 	}
// 	if to, ok := m.(timeouter); ok {
// 		to.SetTimeout(time.Second * 5)
// 	}
// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("%s -> %s", tt.routingKeySend, tt.routingKeyRec), func(t *testing.T) {

// 			ansRecMsg := ""
// 			ansSendMsg := ""
// 			var recTime time.Time
// 			wait := make(chan int)
// 			err = m.Answer(t.Context(), tt.routingKeyRec, func(r string, d amqp.Delivery) ([]byte, error) {
// 				if !strings.EqualFold(r, tt.routingKeySend) {
// 					t.Errorf("Routing keys do not match: have: %s want: %s", r, tt.routingKeySend)
// 				}
// 				ansRecMsg = string(d.Body)
// 				ansSendMsg = fmt.Sprintf("answer-%s", ansRecMsg)
// 				recTime = time.Now()
// 				wait <- 1
// 				return []byte(ansSendMsg), nil
// 			})
// 			if err != nil {
// 				t.Errorf("cannot answer: %v", err)
// 			}
// 			msg := "Test message Ask/Answer WildCard"
// 			sendTime := time.Now()
// 			resp, err := m.Ask(t.Context(), tt.routingKeySend, []byte(msg))
// 			if err != nil {
// 				t.Fatalf("cannot ask: %v", err)
// 			}
// 			<-wait
// 			if !strings.EqualFold(msg, ansRecMsg) {
// 				t.Fatalf("Got message %s expected %s", ansRecMsg, msg)
// 			}
// 			if !strings.EqualFold(ansSendMsg, string(resp.Body)) {
// 				t.Fatalf("Got answer %s expected %s", ansSendMsg, resp.Body)
// 			}
// 			if recTime.Sub(sendTime) > time.Millisecond {
// 				t.Errorf("sending took too long: %s", recTime.Sub(sendTime))
// 			}
// 			close(wait)
// 		}) // synctest
// 	}
// }

package bus

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"log/slog"

	"github.com/vogtp/som/pkg/core/log"
)

type eventer interface {
}

type eventHandler[M eventer] struct {
	wgMsg     sync.WaitGroup
	mu        sync.Mutex
	log       *slog.Logger
	bus       *Manager
	msgType   string
	unsubFucs []unsubscribeFunc
}

func newHandler[M eventer](log *slog.Logger, b *Manager, msgType string) *eventHandler[M] {
	h := &eventHandler[M]{
		log:     log.With("bus", msgType),
		bus:     b,
		msgType: msgType,
	}
	return h
}

// SendSzenarioEvt sends a SzenarioEvtMsg
func (h *eventHandler[M]) Send(evt *M) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.wgMsg.Add(1)
	defer time.AfterFunc(100*time.Millisecond, h.wgMsg.Done)
	b, err := json.Marshal(evt)
	if err != nil {
		h.log.Error("cannot marshal", "event", evt, log.Error, err)
		return fmt.Errorf("cannot marshal %+v: %v", evt, err)
	}
	h.log.Debug("Sending msg", "type", h.msgType, "event", evt)
	h.bus.Emit(h.msgType, b)
	return nil
}

// EventHandler handles events
type EventHandler[M eventer] func(*M)

// HandleSzenarioEvt handles SzenarioEvtMsgs
func (h *eventHandler[M]) Handle(f EventHandler[M]) {
	unsub, err := h.bus.Receive(h.msgType, func(subject string, m *Message) {
		evt := new(M)
		err := json.Unmarshal(m.Body, evt)
		if err != nil {
			h.log.Error("Could not unmarshal message", "payload", string(m.Body), log.Error, err)
		}
		f(evt)
	})
	if err != nil {
		h.log.Error("Cannot handle events", "eventType", h.msgType, log.Error, err)
	}
	h.unsubFucs = append(h.unsubFucs, unsub)
}

func (h *eventHandler[M]) WaitMsgProcessed() {
	h.wgMsg.Wait()
}

func (h *eventHandler[M]) cleanup() {
	h.WaitMsgProcessed()
	for _, us := range h.unsubFucs {
		us()
	}
}

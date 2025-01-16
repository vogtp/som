package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"log/slog"

	"github.com/suborbital/e2core/foundation/bus/bus"
	"github.com/vogtp/som/pkg/core/log"
)

type eventer interface {
}

type eventHandler[M eventer] struct {
	wgMsg    sync.WaitGroup
	mu       sync.Mutex
	log      *slog.Logger
	bus     *bus.Bus
	handlers []*bus.Pod
	msgType  string
}

func newHandler[M eventer](log *slog.Logger, b *bus.Bus, msgType string) *eventHandler[M] {
	h := &eventHandler[M]{
		log:      log.With("bus", msgType),
		bus:     b,
		msgType:  msgType,
		handlers: make([]*bus.Pod, 0),
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
	p := h.bus.Connect()
	defer p.Disconnect()
	p.Send(bus.NewMsg(h.msgType, b))
	return nil
}

// EventHandler handles events
type EventHandler[M eventer] func(*M)

// HandleSzenarioEvt handles SzenarioEvtMsgs
func (h *eventHandler[M]) Handle(f EventHandler[M]) {
	p := h.bus.Connect()
	h.handlers = append(h.handlers, p)
	p.OnType(h.msgType, func(m bus.Message) error {
		h.wgMsg.Add(1)
		defer h.wgMsg.Done()
		evt := new(M)
		err := json.Unmarshal(m.Data(), evt)
		if err != nil {
			h.log.Error("Could not unmarshal message", "payload", string(m.Data()), log.Error, err)
			// does not return an error the the program, just signals the bus
			return fmt.Errorf("could not unmarshal message %s: %w", string(m.Data()), err)
		}
		f(evt)
		return nil
	})
}

func (h *eventHandler[M]) WaitMsgProcessed() {
	h.wgMsg.Wait()
}

func (h *eventHandler[M]) cleanup() {
	h.WaitMsgProcessed()
	for _, h := range h.handlers {
		h.Disconnect()
	}
}

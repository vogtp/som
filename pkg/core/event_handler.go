package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/suborbital/grav/grav"
	"github.com/vogtp/go-hcl"
)

type eventer interface {
}

type eventHandler[M eventer] struct {
	wgMsg    sync.WaitGroup
	mu       sync.Mutex
	hcl      hcl.Logger
	grav     *grav.Grav
	handlers []*grav.Pod
	msgType  string
}

func newHandler[M eventer](hcl hcl.Logger, b *grav.Grav, msgType string) *eventHandler[M] {
	h := &eventHandler[M]{
		hcl:      hcl.Named(msgType),
		grav:     b,
		msgType:  msgType,
		handlers: make([]*grav.Pod, 0),
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
		h.hcl.Errorf("cannot marshal %+v: %v", evt, err)
		return fmt.Errorf("cannot marshal %+v: %v", evt, err)
	}
	h.hcl.Tracef("Sending %s msg %+v", h.msgType, evt)
	p := h.grav.Connect()
	defer p.Disconnect()
	p.Send(grav.NewMsg(h.msgType, b))
	return nil
}

// EventHandler handles events
type EventHandler[M eventer] func(*M)

// HandleSzenarioEvt handles SzenarioEvtMsgs
func (h *eventHandler[M]) Handle(f EventHandler[M]) {
	p := h.grav.Connect()
	h.handlers = append(h.handlers, p)
	p.OnType(h.msgType, func(m grav.Message) error {
		h.wgMsg.Add(1)
		defer h.wgMsg.Done()
		evt := new(M)
		err := json.Unmarshal(m.Data(), evt)
		if err != nil {
			h.hcl.Errorf("Could not unmarshal message %s: %v", string(m.Data()), err)
			// does not return an error the the program, just signals the grav
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

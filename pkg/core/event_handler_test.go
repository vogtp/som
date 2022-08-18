package core

import (
	"testing"
	"time"

	"github.com/vogtp/som/pkg/core/msg"
)

// nolint
func TestHandleMonEvt(t *testing.T) {
	tests := []struct {
		msg *msg.SzenarioEvtMsg
	}{
		{msg: &msg.SzenarioEvtMsg{Name: "test"}},
		{msg: &msg.SzenarioEvtMsg{Name: "dslkjfökladjioru dölkfaj öadlksfu8rwö3o9a öalf3kupoi9"}},
	}
	core, close := New("som-test")
	bus := core.Bus()
	defer close()
	for _, tt := range tests {
		var s *msg.SzenarioEvtMsg
		bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
			if tt.msg.Name != e.Name {
				t.Errorf("got wrong message: %v", e.Name)
			}
			s = e
		})
		if err := bus.Szenario.Send(tt.msg); err != nil {
			t.Fatal(err)
		}

		// for s == nil {
		// 	time.Sleep(10 * time.Millisecond)
		// }
		bus.Szenario.WaitMsgProcessed()
		if s == nil {
			t.Error("no message")
			continue
		}
		if tt.msg.Name != s.Name {
			t.Errorf("got wrong message: %v", s)
		}
	}
}

func TestMsgs(t *testing.T) {
	mm := &msg.SzenarioEvtMsg{Name: "BaseMsg", Time: time.Now()}

	if mm.Name != "BaseMsg" {
		t.Errorf("Wron name: %s", mm.Name)
	}
}

func TestTrans(t *testing.T) {
	bus, close := New("som-test")
	defer close()
	_ = bus
}

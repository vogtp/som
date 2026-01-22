package core

import (
	"reflect"
	"sync"
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
	bus := core.BusFIXME()
	defer close()

	rec := make(map[string]msg.SzenarioEvtMsg)
	snd := make(map[string]msg.SzenarioEvtMsg)

	var wg sync.WaitGroup
	wg.Add(len(tests))
	bus.Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
		rec[e.Name] = *e
		wg.Done()
	})

	bus.Alert.Handle(func(am *msg.AlertMsg) {
		t.Errorf("Got an alert message: %v",am.Name)
	})

	for _, tt := range tests {
		if err := bus.Szenario.Send(tt.msg); err != nil {
			t.Fatalf("cannot send msg: %v", err)
		}
		snd[tt.msg.Name] = *tt.msg
	}
	wg.Wait()
	if !reflect.DeepEqual(snd, rec) {
		t.Errorf("Send and receive not equal.\nrec: %+v\nsnd: %+v", rec, snd)
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

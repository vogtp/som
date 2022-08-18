package alertmgr

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
)

func TestAlertMgr_checkEvent(t *testing.T) {
	core.New("")
	am := AlertMgr{hcl: hcl.New(), basicStates: make(map[string]*basicState), status: status.New()}
	am.alertIntervall = time.Hour
	viper.Set(cfg.AlertDelay, time.Microsecond)
	evtA := msg.NewSzenarioEvtMsg("A", "user", time.Now())
	evtB := msg.NewSzenarioEvtMsg("B", "user", time.Now())

	if am.checkEvent(evtA) != nil {
		t.Error("Should not generate alert")
	}
	if am.checkEvent(evtA) != nil {
		t.Error("Should not generate alert")
	}
	if am.checkEvent(evtB) != nil {
		t.Error("Should not generate alert")
	}
	evtA.AddErr(errors.New("test"))
	if am.checkEvent(evtA) != nil {
		t.Error("Should not generate alert")
	}
	evtA.AddErr(errors.New("test2"))
	if am.checkEvent(evtA) == nil {
		t.Error("Should  generate alert")
	}
	if am.checkEvent(evtB) != nil {
		t.Error("Should not generate alert")
	}
	if am.checkEvent(evtA) != nil {
		t.Error("Should not generate alert")
	}
	evtA1 := msg.NewSzenarioEvtMsg("A", "user", time.Now().Add(-25*time.Hour))
	evtA1.AddErr(errors.New("test"))
	if am.checkEvent(evtA1) != nil {
		t.Error("Should generate alert")
	}
}

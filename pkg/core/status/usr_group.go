package status

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
)

// UserGroup correlates event messages
type UserGroup interface {
	EvtGroup
	slog.LogValuer
	String() string
	StringInt(int) string
	AddEvent(evt *msg.SzenarioEvtMsg)
	LastEvent() *msg.SzenarioEvtMsg
	LastUpdate() time.Time
	LastTotal() float64
	Totals() []float64
	Error() string
}

type usrGroup struct {
	*evtGroup
}

func (usrGroup) New(k string) Grouper {
	return &usrGroup{&evtGroup{
		maxLen: -1,
		Group:  &Group{key: k, children: make([]Grouper, 0)}},
	}
}

// AddEvent to be correlated
func (ug *usrGroup) AddEvent(evt *msg.SzenarioEvtMsg) {
	ug.evtGroup.AddEvent(evt)
}

func (ug usrGroup) String() string {
	return ug.StringInt(0)
}

func (ug usrGroup) StringInt(i int) string {
	it := ""
	for c := 0; c < i; c++ {
		it += " "
	}
	str := fmt.Sprintf("%sUser %s: %s", it, ug.Key(), ug.Level())
	cld := ug.children
	for i := len(cld) - 1; i >= 0; i-- {
		c := cld[i]
		if e, ok := c.(*evtWrapper); ok {
			stat := e.Level().String()
			if e.evt.Err() != nil {
				stat = e.evt.Err().Error()
			}
			t := getEvtTot(e.evt)
			if t > 0 {
				stat = fmt.Sprintf("%5.2fs %s", t, stat)
			}

			str = fmt.Sprintf("%s\n%s%s%s: %s", str, it, it, e.evt.Time.Format(cfg.TimeFormatString), stat)
		} else {
			str = fmt.Sprintf("%s\n%s%s%s: %s", str, it, it, c.Level(), c.Key())
		}
	}
	return str
}

// LogValue satisfies slog
func (ug usrGroup) LogValue() slog.Value {
	rv := make([]slog.Attr, 0, 3)
	rv = append(rv, slog.String(log.Szenario, ug.Key()))
	rv = append(rv, slog.String("alert_level", ug.Level().String()))
	cld := ug.children
	hist := make([]slog.Attr, 0, len(ug.children))
	for i := len(cld) - 1; i >= 0; i-- {
		c := cld[i]
		if e, ok := c.(*evtWrapper); ok {
			stat := e.Level().String()
			if e.evt.Err() != nil {
				stat = e.evt.Err().Error()
			}
			t := getEvtTot(e.evt)
			hist = append(hist, slog.Group(fmt.Sprintf("%v",e.evt.Time.Unix()), slog.Time(slog.TimeKey, e.evt.Time), slog.Float64("duration", t), slog.String("status", stat)))
		}
	}
	rv = append(rv, slog.Any("history", slog.GroupValue(hist...)))
	return slog.GroupValue(rv...)
}

func (ug usrGroup) LastUpdate() time.Time {
	var lu time.Time
	for _, c := range ug.children {
		if eg, ok := c.(*evtWrapper); ok {
			if eg.evt.Time.After(lu) {
				lu = eg.evt.Time
			}
		}
	}
	return lu
}

func (ug usrGroup) LastEvent() *msg.SzenarioEvtMsg {
	l := len(ug.children)
	if l < 1 {
		return nil
	}
	c := ug.children[l-1]
	if e, ok := c.(*evtWrapper); ok {
		return e.evt
	}
	return nil
}

func getEvtTot(e *msg.SzenarioEvtMsg) float64 {
	if tot, ok := e.Counters["step.total"]; ok {
		return tot
	}
	return 0
}

func (ug usrGroup) LastTotal() float64 {
	e := ug.LastEvent()
	if e == nil || e.Err() != nil {
		return 0
	}
	return getEvtTot(e)
}

func (ug usrGroup) Totals() []float64 {
	tot := make([]float64, 0)
	for _, c := range ug.children {
		if e, ok := c.(*evtWrapper); ok {
			t := getEvtTot(e.evt)
			if t > 0 && e.evt.Err() == nil {
				tot = append(tot, t)
			}
		}
	}
	return tot
}

func (ug usrGroup) Error() string {
	e := ug.LastEvent()
	if e != nil && e.Err() != nil {
		return e.Err().Error()
	}
	return ""
}

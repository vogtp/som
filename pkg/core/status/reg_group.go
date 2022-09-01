package status

import (
	"fmt"
	"time"

	"github.com/vogtp/som/pkg/core/msg"
)

// RegionGroup correlates event messages
type RegionGroup interface {
	Grouper
	String() string
	StringInt(int) string
	AddEvent(evt *msg.SzenarioEvtMsg)
	Users() []UserGroup
	LastUpdate() time.Time
	LastTotal() float64
	Totals() []float64
}

type regGroup struct {
	*Group
}

func (regGroup) New(k string) Grouper {
	return &regGroup{Group: &Group{key: k, children: make([]Grouper, 0)}}
}

// AddEvent to be correlated
func (rg *regGroup) AddEvent(evt *msg.SzenarioEvtMsg) {
	rg.getOrCreateGroup(evt.Username).AddEvent(evt)
}

func (rg *regGroup) getOrCreateGroup(key string) UserGroup {
	for _, c := range rg.children {
		if c.Key() == key {
			return c.(UserGroup)
		}
	}
	c := &usrGroup{
		evtGroup: &evtGroup{
			Group: &Group{
				cfg:      rg.cfg,
				key:      key,
				children: make([]Grouper, 0),
			},
			maxLen: -1,
		},
	}
	rg.Add(c)
	return c
}

func (rg regGroup) String() string {
	return rg.StringInt(0)
}

func (rg regGroup) StringInt(i int) string {
	it := ""
	for c := 0; c < i; c++ {
		it += " "
	}
	str := fmt.Sprintf("%sRegion %s: %s", it, rg.Key(), rg.Level())
	for _, c := range rg.Users() {
		str = fmt.Sprintf("%s\n%s%s", str, it, c.StringInt(i+1))
	}
	return str
}

func (rg regGroup) Users() []UserGroup {
	rgs := make([]UserGroup, len(rg.children))
	for i, c := range rg.children {
		rgs[i] = c.(UserGroup)
	}
	return rgs
}

func (rg regGroup) LastUpdate() time.Time {
	var lu time.Time
	for _, c := range rg.Users() {
		t := c.LastUpdate()
		if t.After(lu) {
			lu = t
		}
	}
	return lu
}

func (rg regGroup) LastTotal() float64 {
	t := 0.
	n := 0.
	for _, c := range rg.Users() {
		v := c.LastTotal()
		if v == 0 {
			continue
		}
		t += v
		n++
	}
	return t / n
}

func (rg regGroup) Totals() []float64 {
	tot := make([]float64, 0)
	for _, c := range rg.Users() {
		tot = append(tot, c.Totals()...)
	}
	return tot
}

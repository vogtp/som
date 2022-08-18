package status

import (
	"fmt"
	"time"

	"github.com/vogtp/som/pkg/core/msg"
)

// SzenarioGroup correlates event messages
type SzenarioGroup interface {
	Grouper
	String() string
	StringInt(int) string
	AddEvent(evt *msg.SzenarioEvtMsg)
	Regions() []RegionGroup
	LastUpdate() time.Time
	LastTotal() float64
	Totals() []float64
}

type szGroup struct {
	*Group
}

func (szGroup) New(k string) Grouper {
	return &szGroup{Group: &Group{key: k, children: make([]Grouper, 0)}}
}

// AddEvent to be correlated
func (sg *szGroup) AddEvent(evt *msg.SzenarioEvtMsg) {
	sg.getOrCreateGroup(evt.Region).AddEvent(evt)
}

func (sg *szGroup) getOrCreateGroup(key string) RegionGroup {
	for _, c := range sg.children {
		if c.Key() == key {
			return c.(RegionGroup)
		}
	}
	c := &regGroup{
		Group: &Group{
			key:      key,
			children: make([]Grouper, 0),
		},
	}
	sg.Add(c)
	return c
}

func (sg szGroup) String() string {
	return sg.StringInt(-1) // there is no output on this level -> do not intent the next
}

func (sg szGroup) StringInt(i int) string {
	it := ""
	for c := 0; c < i; c++ {
		it += " "
	}
	str := fmt.Sprintf("%sSzenario %s: %s", it, sg.Key(), sg.Level())
	for _, c := range sg.Regions() {
		str = fmt.Sprintf("%s\n%s%s", str, it, c.StringInt(i+1))
	}
	return str
}

func (sg szGroup) Regions() []RegionGroup {
	rgs := make([]RegionGroup, len(sg.children))
	for i, c := range sg.children {
		rgs[i] = c.(RegionGroup)
	}
	return rgs
}

func (sg szGroup) LastUpdate() time.Time {
	var lu time.Time
	for _, c := range sg.Regions() {
		t := c.LastUpdate()
		if t.After(lu) {
			lu = t
		}
	}
	return lu
}

func (sg szGroup) LastTotal() float64 {
	t := 0.
	n := 0.
	for _, c := range sg.Regions() {
		v := c.LastTotal()
		if v == 0 {
			continue
		}
		t += v
		n++
	}
	return t / n
}

func (sg szGroup) Totals() []float64 {
	tot := make([]float64, 0)
	for _, c := range sg.Regions() {
		tot = append(tot, c.Totals()...)
	}
	return tot
}

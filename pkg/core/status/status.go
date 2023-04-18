package status

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/msg"
)

// Status correlates event messages
type Status interface {
	Grouper
	String() string
	StringInt(int) string
	AddEvent(evt *msg.SzenarioEvtMsg)
	Szenarios() []SzenarioGroup
	Get(string) SzenarioGroup
	LastUpdate() time.Time
	JSONBySzenario(string) []byte
	UpdatePrometheus()
}

type statusGroup struct {
	*Group
	promAvail *promAvail
}

// New Creates a TopGroup to correlate event messages
func New() Status {
	sg := &statusGroup{
		Group: &Group{
			key:      "root",
			children: make([]Grouper, 0),
		},
	}
	return sg
}

func (statusGroup) New(k string) Grouper {
	return &statusGroup{Group: &Group{key: k, children: make([]Grouper, 0)}}
}

// AddEvent to be correlated
func (sg *statusGroup) AddEvent(evt *msg.SzenarioEvtMsg) {
	sg.getOrCreateGroup(evt.Name).AddEvent(evt)
}

func (sg *statusGroup) getOrCreateGroup(key string) SzenarioGroup {
	for _, c := range sg.children {
		if c.Key() == key {
			return c.(SzenarioGroup)
		}
	}
	c := &szGroup{
		Group: &Group{
			cfg:      sg.cfg,
			key:      key,
			children: make([]Grouper, 0),
		},
	}
	sg.Add(c)
	return c
}

func (sg statusGroup) String() string {
	return sg.StringInt(0)
}

func (sg statusGroup) StringInt(i int) string {
	it := ""
	for c := 0; c < i; c++ {
		it += " "
	}
	str := "" //fmt.Sprintf("%s%s: %s", it, "Root", sg.Level())
	for _, c := range sg.Szenarios() {
		str = fmt.Sprintf("%s\n%s%s", str, it, c.StringInt(i+1))
	}
	return str
}

func (sg statusGroup) Get(key string) SzenarioGroup {
	for _, s := range sg.children {
		if s.Key() == key {
			return s.(SzenarioGroup)
		}
	}
	return nil
}

func (sg statusGroup) Szenarios() []SzenarioGroup {
	rgs := make([]SzenarioGroup, len(sg.children))
	for i, c := range sg.children {
		rgs[i] = c.(SzenarioGroup)
	}
	sort.Slice(rgs, func(i, j int) bool {
		return strings.ToLower(rgs[i].Key()) < strings.ToLower(rgs[j].Key())
	})
	return rgs
}

func (sg statusGroup) LastUpdate() time.Time {
	var lu time.Time
	for _, c := range sg.Szenarios() {
		t := c.LastUpdate()
		if t.After(lu) {
			lu = t
		}
	}
	return lu
}

func (sg statusGroup) JSONBySzenario(n string) []byte {
	g := New().(*statusGroup)
	g.key = n
	for _, sz := range sg.children {
		if sz.Key() != n {
			continue
		}
		g.children = append(g.children, sz)
		break
	}
	d, err := g.MarshalJSON()
	if err != nil {
		hcl.Warn("Cannot marshal JSONBySzenario", "key", n, "error", err)
	}
	return d
}

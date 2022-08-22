package status

import (
	"encoding/json"
	"fmt"
	"math"
)

// Availability is the availability of a group
// given in % from 0 to 1
type Availability float64

func (ga Availability) String() string {
	return fmt.Sprintf("%3.0f%%", ga*100)
}

// Grouper is the alert group interface
type Grouper interface {
	New(k string) Grouper
	Key() string
	Level() Level
	Availability() Availability
	Add(Grouper)
	json.Marshaler
	json.Unmarshaler
}

// Group is the implementation of a generic alert group
type Group struct {
	children []Grouper
	key      string
}

type jsonGroup struct {
	Key      string
	Children []jsonChild
}

type jsonChild struct {
	Typ  string
	Key  string
	JSON string
}

// New create a Group
func (Group) New(k string) Grouper {
	return &Group{key: k, children: make([]Grouper, 0)}
}

// Add a new child
func (g *Group) Add(c Grouper) {
	if g.children == nil {
		g.children = make([]Grouper, 0)
	}
	g.children = append(g.children, c)
}

// Key returns the key of the group item
func (g Group) Key() string {
	return g.key
}

// Availability is the availability of a group
// given in % from 0 to 1
func (g Group) Availability() Availability {
	a := Availability(0.0)
	for _, c := range g.children {
		a += c.Availability()
	}
	a /= Availability(len(g.children))
	return a
}

// Level returns the error level
func (g Group) Level() Level {
	return g.LevelRedundent()
}

// LevelRedundent returns the error level of a redanent set
func (g Group) LevelRedundent() Level {
	childCnt := float64(len(g.children))
	if childCnt < 1 {
		return Unknown
	}
	var tot float64
	for _, c := range g.children {
		tot += float64(c.Level())
	}
	tot /= childCnt
	// handle Down strictly
	if tot > float64(Warning) && tot < (float64(Down)-.2) {
		tot = float64(Warning)
	}
	// only have issues when issue level is reached
	if tot < float64(Issues) && tot > float64(OK) {
		tot = float64(OK)
	}
	// only have unknown when unknown level is reached
	if tot < float64(OK) && tot > float64(Unknown) {
		tot = float64(OK)
	}
	lvl := math.Round(tot)
	return Level(lvl)
}

// MarshalJSON marshall children
func (g Group) MarshalJSON() ([]byte, error) {
	j := jsonGroup{
		Key:      g.key,
		Children: make([]jsonChild, len(g.children)),
	}
	for i, c := range g.children {
		d, err := c.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("cannot marshal child: %v", err)
		}
		j.Children[i] = jsonChild{
			Typ:  fmt.Sprintf("%T", c),
			Key:  c.Key(),
			JSON: string(d),
		}
	}

	return json.Marshal(&j)
}

// UnmarshalJSON marshall children
func (g *Group) UnmarshalJSON(data []byte) error {
	j := &jsonGroup{
		Children: make([]jsonChild, 0),
	}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	g.key = j.Key
	for _, jChild := range j.Children {
		var c Grouper
		c, err := GrpReg.new(jChild.Typ, jChild.Key)
		if err != nil {
			return fmt.Errorf("cannot unmarshall group: %v: %v", jChild.Typ, err)
		}
		if err := c.UnmarshalJSON([]byte(jChild.JSON)); err != nil {
			return fmt.Errorf("cannot unmarshal child: %v", err)
		}
		g.Add(c)

	}
	return nil
}

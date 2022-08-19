package status

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
)

type testLevel struct {
	l Level
	k string
}

func (testLevel) New(k string) Grouper { return &testLevel{k: k} }
func (l *testLevel) Key() string {
	if len(l.k) < 1 {
		l.k = uuid.NewString()
	}
	return l.k
}
func (l testLevel) Add(Grouper)                  { panic("operation not allowed") }
func (l testLevel) Level() Level                 { return Level(l.l) }
func (l testLevel) MarshalJSON() ([]byte, error) { return []byte(Level(l.l).String()), nil }
func (l testLevel) Availability() Groupavailability {
	if l.l == OK {
		return 1
	}
	return 0
}

func (l *testLevel) UnmarshalJSON(data []byte) error {
	s := string(data)
	for i := Unknown; i <= Down; i++ {
		if i.String() == s {
			l.l = i
			return nil
		}
	}
	return fmt.Errorf("level could not find level: %s", s)
}

func init() {
	GrpReg.Add(ok)
	viper.SetDefault(cfg.AlertIncidentCorrelationEvents, 3)
}

var (
	down    *testLevel          = &testLevel{l: Down}
	warn    *testLevel          = &testLevel{l: Warning}
	issues  *testLevel          = &testLevel{l: Issues}
	ok      *testLevel          = &testLevel{l: OK}
	unknown *testLevel          = &testLevel{l: Unknown}
	evtOK   *msg.SzenarioEvtMsg = &msg.SzenarioEvtMsg{Time: time.Now()}
	evtErr  *msg.SzenarioEvtMsg = &msg.SzenarioEvtMsg{Time: time.Now()}
)

func testLevelAndJSONEnc(t *testing.T, g Grouper, lvl Level) {
	t.Helper()
	if got := g.Level(); got != lvl {
		t.Errorf("Group.Level() = %v, want %v", got, lvl)
	}

	data, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("Cannot marshal group: %v", err)
	}
	g2 := &Group{}
	err = json.Unmarshal(data, g2)
	if err != nil {
		t.Errorf("cannot unmarshal group: %v -- %s", err, string(data))
		return
	}
	if got := g2.Level(); got != g.Level() {
		t.Errorf("json converted Group.Level() = %v, want %v -- %s", got, g.Level(), string(data))
	}
}

func TestGroup_JSON(t *testing.T) {
	okGrp := &Group{}
	okGrp.Add(ok)
	downGrp := &Group{}
	downGrp.Add(down)
	warnGrp := &Group{}
	warnGrp.Add(downGrp)
	warnGrp.Add(okGrp)

	evtGrpOK := newEvtGroup("test")
	evtGrpOK.AddEvent(evtErr)
	evtGrpOK.AddEvent(evtErr)
	evtGrpOK.AddEvent(evtErr)
	evtGrpOK.AddEvent(evtOK)

	tests := []struct {
		name     string
		children []Grouper
		want     Level
	}{
		{"dow", []Grouper{downGrp, okGrp, warnGrp}, Warning},
		{"dowo", []Grouper{downGrp, okGrp, warnGrp, okGrp}, Issues},
		{"eO", []Grouper{evtGrpOK}, Unknown},
		{"eOw", []Grouper{evtGrpOK, warnGrp}, OK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{key: tt.name}
			for _, l := range tt.children {
				g.Add(l)
			}

			if tt.name != g.Key() {
				t.Errorf("wrong key: %v", g.Key())
			}
			testLevelAndJSONEnc(t, g, tt.want)
		})
	}
}

func TestGroup_Level(t *testing.T) {
	tests := []struct {
		name     string
		children []*testLevel
		want     Level
	}{
		{"allOK", []*testLevel{ok, ok, ok, ok, ok}, OK},
		{"allWarn", []*testLevel{warn, warn, warn}, Warning},
		{"allDown", []*testLevel{down, down, down}, Down},
		{"allIssues", []*testLevel{issues, issues}, Issues},
		{"allUnknow", []*testLevel{unknown, unknown, unknown}, Unknown},
		{"oneDown", []*testLevel{down}, Down},
		{"nothing", []*testLevel{}, Unknown},

		{"4OK1down", []*testLevel{ok, ok, ok, ok, down}, OK},
		{"3OK2down", []*testLevel{ok, ok, ok, down, down}, Issues},
		{"2OK1warn2down", []*testLevel{ok, ok, warn, down, down}, Warning},
		{"2OK2down", []*testLevel{ok, ok, down, down}, Warning},
		{"2warn2down", []*testLevel{warn, warn, down, down}, Warning},
		{"2warn1iss1down", []*testLevel{warn, warn, issues, down}, Warning},
		{"1OK3down", []*testLevel{ok, down, down, down}, Warning},
		{"2warn2Issues", []*testLevel{warn, warn, issues, issues}, Warning},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{key: tt.name}
			for _, l := range tt.children {
				g.Add(l)
			}

			if tt.name != g.Key() {
				t.Errorf("wrong key: %v", g.Key())
			}
			testLevelAndJSONEnc(t, g, tt.want)
		})
	}
}

func TestGroup_LevelHyrachical(t *testing.T) {
	downGrp := &Group{key: "downGrp"}
	downGrp.Add(down)
	downGrp.Add(down)
	if downGrp.Level() != Down {
		t.Fatalf("downGrp should have level down not %v", downGrp.Level())
	}
	warnGrp := &Group{key: "warnGrp"}
	warnGrp.Add(down)
	warnGrp.Add(down)
	if warnGrp.Level() != Down {
		t.Fatalf("warnGrp should have level warn not %v", warnGrp.Level())
	}
	issueGrp := &Group{key: "issueGrp"}
	issueGrp.Add(issues)
	issueGrp.Add(ok)
	issueGrp.Add(warnGrp)
	if issueGrp.Level() != Issues {
		t.Fatalf("issueGrp should have level issue not %v", issueGrp.Level())
	}
	okGrp := &Group{key: "okGrp"}
	okGrp.Add(ok)
	okGrp.Add(ok)
	okGrp.Add(ok)
	if okGrp.Level() != OK {
		t.Fatalf("okGrp should have level OK not %v", okGrp.Level())
	}
	warnHyraGrp := &Group{key: "warnHyraGrp"}
	warnHyraGrp.Add(down)
	warnHyraGrp.Add(warnGrp)
	warnHyraGrp.Add(ok)
	if warnHyraGrp.Level() != Warning {
		t.Fatalf("warnHyraGrp should have level Warn not %v", warnHyraGrp.Level())
	}
	downHyraGrp := &Group{key: "downHyraGrp"}
	downHyraGrp.Add(down)
	downHyraGrp.Add(warnGrp)
	downHyraGrp.Add(down)
	if downHyraGrp.Level() != Down {
		t.Fatalf("downHyraGrp should have level Down not %v", downHyraGrp.Level())
	}
	testLevelAndJSONEnc(t, downHyraGrp, Down)
}

func TestEvtGroup_Level(t *testing.T) {
	tests := []struct {
		name     string
		children []*testLevel
		want     Level
	}{
		{"allOK", []*testLevel{ok, ok, ok, ok, ok}, OK},
		{"allWarn", []*testLevel{warn, warn, warn}, Warning},
		{"allDown", []*testLevel{down, down, down}, Down},
		{"allIssues", []*testLevel{issues, issues}, Issues},
		{"allUnknow", []*testLevel{unknown, unknown, unknown}, Unknown},
		{"oneDown", []*testLevel{down}, Down},
		{"nothing", []*testLevel{}, Unknown},

		{"4OK1down", []*testLevel{ok, ok, ok, ok, down}, OK},
		{"3OK2down", []*testLevel{ok, ok, ok, down, down}, Issues},
		{"2OK1warn2down", []*testLevel{ok, ok, warn, down, down}, Warning},
		{"2OK2down", []*testLevel{ok, ok, down, down}, Warning},
		{"2warn2down", []*testLevel{warn, warn, down, down}, Warning},
		{"2warn1iss1down", []*testLevel{warn, warn, issues, down}, Warning},
		{"1OK3down", []*testLevel{ok, down, down, down}, Warning},
		{"2warn2Issues", []*testLevel{warn, warn, issues, issues}, Warning},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &evtGroup{Group: &Group{key: tt.name}}
			for _, l := range tt.children {
				g.Add(l)
			}

			if got := g.Level(); got != tt.want {
				t.Errorf("Group.Level() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvtGroup_QueueLenth(t *testing.T) {
	tests := []struct {
		name  string
		child *testLevel
		want  Level
	}{
		{"ok", ok, OK},
		{"okWarn", warn, Issues},
		{"okWarnWarn", warn, Issues},
		{"WarnWarnWarn", warn, Warning},
		{"WarnWarnIss", issues, Warning},
		{"WarnIssIss", issues, Issues},
		{"IssIssDown", down, Warning},
		{"IssDownown", down, Warning},
		{"DownDownown", down, Down},
	}
	g := &evtGroup{maxLen: 3, Group: &Group{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Group.key = tt.name
			g.Add(tt.child)

			if got := g.Level(); got != tt.want {
				t.Errorf("Group.Level() = %v, want %v (%v)", got, tt.want, g.children)
			}
		})
	}
}

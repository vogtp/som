package status

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vogtp/som/pkg/core/msg"
)

func checkSzGrp(t *testing.T, g Grouper, lvl Level) {
	t.Helper()
	if got := g.Level(); got != lvl {
		details := ""
		if str, ok := g.(fmt.Stringer); ok {
			details = fmt.Sprintf("\n%s", str.String())
		}
		t.Errorf("Group.Level() = %v, want %v%s", got, lvl, details)
	}

	data, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("Cannot marshal group: %v", err)
	}
	g2 := New()
	err = json.Unmarshal(data, g2)
	if err != nil {
		t.Errorf("cannot unmarshal group: %v -- %s", err, string(data))
		return
	}
	if got := g2.Level(); got != g.Level() {
		t.Errorf("json converted Group.Level() = %v, want %v -- %s", got, g.Level(), string(data))
	}
}

func Test_szGroup_JSON(t *testing.T) {

	rug := New()

	rug.AddEvent(&msg.SzenarioEvtMsg{
		ID:       uuid.NewString(),
		Name:     "testSzenario1",
		Region:   "testRegion1",
		Username: "testUser1",
	})
	rug.AddEvent(&msg.SzenarioEvtMsg{
		ID:       uuid.NewString(),
		Name:     "testSzenario2",
		Region:   "testRegion1",
		Username: "testUser1",
	})
	b, err := json.Marshal(rug)
	if err != nil {
		t.Fatalf("cannot marshal: %v", err)
	}
	rug2 := New()
	err = json.Unmarshal(b, rug2)
	if err != nil {
		t.Fatalf("cannot unmarshal: %v", err)
	}
	deep.MaxDepth = 100
	if diff := deep.Equal(rug, rug2); diff != nil {
		t.Error(diff)
	}
	// if !reflect.DeepEqual(rug, rug2) {
	// 	t.Errorf("Not deep equal:\n%v", string(b))
	// }
}

func Test_szGroup_AddEvent(t *testing.T) {
	rug := New()
	rugO, ok := rug.(*statusGroup)
	if !ok {
		t.Fatal("Cannot cast rug")
	}
	rug.AddEvent(&msg.SzenarioEvtMsg{
		ID:       uuid.NewString(),
		Name:     "testSzenario1",
		Region:   "testRegion1",
		Username: "testUser1",
	})
	checkSzGrp(t, rug, OK)

	evt := &msg.SzenarioEvtMsg{
		ID:       uuid.NewString(),
		Name:     "testSzenario1",
		Region:   "testRegion1",
		Username: "testUser1",
	}
	evt.AddErr(errors.New("test error"))
	rug.AddEvent(evt)
	checkSzGrp(t, rug, Warning)
	if len(rugO.children) > 1 {
		t.Errorf("Should only have one reqion not %v", len(rugO.children))
	}
	evt = &msg.SzenarioEvtMsg{
		ID:       uuid.NewString(),
		Name:     "testSzenario2",
		Region:   "testRegion1",
		Username: "testUser1",
	}
	evt.AddErr(errors.New("test error"))
	rug.AddEvent(evt)
	checkSzGrp(t, rug, Warning)
	evt = &msg.SzenarioEvtMsg{
		ID:       uuid.NewString(),
		Name:     "testSzenario1",
		Region:   "testRegion1",
		Username: "testUser1",
	}
	rug.AddEvent(evt)
	checkSzGrp(t, rug, Warning)
	evt = &msg.SzenarioEvtMsg{
		ID:       uuid.NewString(),
		Name:     "testSzenario2",
		Region:   "testRegion1",
		Username: "testUser1",
	}
	rug.AddEvent(evt)
	checkSzGrp(t, rug, OK)
}

package webstatus

import (
	"testing"
	"time"
)

func Test_dayTickMarks_Ticks(t *testing.T) {
	tests := []struct {
		numTicks int
	}{
		{1}, {3}, {7},
	}
	for _, tt := range tests {
		t1 := time.Now()
		t2 := t1.Add(time.Duration(tt.numTicks) * 24 * time.Hour)

		var ticks dayTickMarks

		res := ticks.Ticks(float64(t1.Unix()), float64(t2.Unix()))
		if len(res) > tt.numTicks*48+1 {
			t.Errorf("Should not have more than %v ticks have %v", tt.numTicks*48, len(res))
		}
	}
}

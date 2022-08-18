package webstatus

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

func TestWebStatus_getIncidentDetailFiles(t *testing.T) {
	_, close := core.New("som.visualiser-test")
	defer close()
	s := New()
	viper.Set(cfg.DataDir, "../../../data/")
	start := time.Now()
	files, err := s.getIncidentFiles(s.getIncidentRoot(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 1 {
		t.Fatal("No files returned")
	}
	d := time.Since(start)
	d /= time.Duration(len(files))
	if d > 40*time.Millisecond {
		t.Errorf("Incident loading took too long: %v", d)
	}
}

package cdp

import (
	"time"

	"github.com/iancoleman/strcase"
	"github.com/vogtp/go-hcl"
)

type stepInfo struct {
	hcl       *hcl.Logger
	name      string
	startTime time.Time
	stepTimes map[string]float64
}

func newStepInfo(hcl *hcl.Logger) *stepInfo {
	return &stepInfo{
		hcl:       hcl,
		stepTimes: make(map[string]float64),
	}
}

func (s *stepInfo) start(name string) {
	s.hcl.Debugf("Step: %s", name)
	s.name = name
	s.startTime = time.Now()
}

func (s *stepInfo) end(name string) {
	if s.startTime.IsZero() {
		panic("step start time must not be zero")
	}
	d := time.Since(s.startTime)
	if d.Seconds() > .2 {
		s.stepTimes[strcase.ToCamel(s.name)] = d.Seconds()
	}
	s.hcl.Infof("Step %q took %v", s.name, d)
}

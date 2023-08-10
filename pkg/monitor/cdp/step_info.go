package cdp

import (
	"time"

	"log/slog"

	"github.com/iancoleman/strcase"
)

type stepInfo struct {
	log       *slog.Logger
	name      string
	startTime time.Time
	stepTimes map[string]float64
}

func newStepInfo(log *slog.Logger) *stepInfo {
	return &stepInfo{
		log:       log,
		stepTimes: make(map[string]float64),
	}
}

func (s *stepInfo) start(name string) {
	s.log.Debug("Step start", "step", name)
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
	s.log.Info("Step finished", "step", s.name, "duration", d)
}

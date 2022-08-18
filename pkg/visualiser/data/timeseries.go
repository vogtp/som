package data

import "time"

type Timeserie struct {
	Data []DataPoint
}

func (s *Timeserie) Add(t time.Time, p float64) {
	if s.Data == nil {
		s.Data = make([]DataPoint, 0)
	}
	s.Data = append([]DataPoint{{t, p}}, s.Data...)
}

func (s Timeserie) Avg() float64 {
	if len(s.Data) < 1 {
		return 0
	}
	var t float64
	for _, p := range s.Data {
		t += p.Point
	}
	return t / float64(len(s.Data))
}

func (s Timeserie) Len() int {
	return len(s.Data)
}

func (s Timeserie) GetSlices() ([]time.Time, []float64) {

	t := make([]time.Time, s.Len())
	p := make([]float64, s.Len())

	for i, d := range s.Data {
		t[i] = d.Time
		p[i] = d.Point
	}
	return t, p
}

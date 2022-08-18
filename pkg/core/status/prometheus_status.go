package status

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vogtp/som/pkg/bridger"
)

type promAvail struct {
	availGauges map[string]prometheus.Gauge
}

// UpdatePrometheus updates the prometheus counters
func (sg *statusGroup) UpdatePrometheus() {
	if !bridger.PrometheusIsActive() {
		return
	}
	if sg.promAvail == nil {
		sg.promAvail = &promAvail{availGauges: make(map[string]prometheus.Gauge)}
	}
	for _, sz := range sg.Szenarios() {
		sg.promAvail.updateSzenario(sz)
	}
}

func (pa *promAvail) updateSzenario(sz SzenarioGroup) {
	pa.getAvailGauge(sz).Set(float64(sz.Availability()))
}

func (pa *promAvail) getAvailGauge(sz SzenarioGroup) prometheus.Gauge {
	name := fmt.Sprintf("%s_availability", bridger.PrometheusName(sz.Key()))
	ag := pa.availGauges[name]
	if ag == nil {
		ag = promauto.NewGauge(prometheus.GaugeOpts{Name: name})
		pa.availGauges[name] = ag
	}
	return ag
}

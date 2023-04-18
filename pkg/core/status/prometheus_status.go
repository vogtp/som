package status

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/bridger"
	"github.com/vogtp/som/pkg/core/log"
)

type promAvail struct {
	availGauges    map[string]prometheus.Gauge
	availGaugeVecs map[string]*prometheus.GaugeVec
}

// UpdatePrometheus updates the prometheus counters
func (sg *statusGroup) UpdatePrometheus() {
	if !bridger.PrometheusIsActive() {
		return
	}
	if sg.promAvail == nil {
		sg.promAvail = &promAvail{
			availGauges:    make(map[string]prometheus.Gauge),
			availGaugeVecs: make(map[string]*prometheus.GaugeVec),
		}
	}
	for _, sz := range sg.Szenarios() {
		sg.promAvail.updateSzenario(sz)
	}
}

func (pa *promAvail) updateSzenario(sz SzenarioGroup) {
	avail := float64(sz.Availability())
	pa.getAvailGauge(sz).Set(avail)
	gv := pa.getAvailGaugeVec(sz)
	gvAvail, err := gv.GetMetricWithLabelValues("availability")
	if err != nil {
		hcl.Warn("Cannot get gauge vec", log.Error, err, "counter", "availability", "key", sz.Key())
		return
	}
	gvAvail.Set(avail * 100)
	gvTotal, err := gv.GetMetricWithLabelValues("check_time")
	if err != nil {
		hcl.Warn("Cannot get gauge vec: %v", log.Error, err, "counter", "check_time", "key", sz.Key())
		return
	}
	gvTotal.Set(sz.LastTotal())
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

func (pa *promAvail) getAvailGaugeVec(sz SzenarioGroup) *prometheus.GaugeVec {
	name := fmt.Sprintf("%s_summary", bridger.PrometheusName(sz.Key()))
	gv := pa.availGaugeVecs[name]
	if gv == nil {
		gv = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: name,
		},
			[]string{"type"},
		)
		pa.availGaugeVecs[name] = gv
	}
	return gv
}

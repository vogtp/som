package bridger

import (
	"errors"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

const (
	downValue = 0
)

// RegisterPrometheus registers NetCrunch Messages on the eventbus
func RegisterPrometheus() {
	bus := core.Get().Bus()
	p := prometheusBackend{
		hcl:      bus.GetLogger().Named("prometheus"),
		promSz:   make(map[string]*promSz),
		gaugeVec: make(map[string]*prometheus.GaugeVec),
		histoVec: make(map[string]*prometheus.HistogramVec),
	}
	bus.Szenario.Handle(p.handleEventBus)
	go p.start()
}

var (
	active bool = false
)

// PrometheusIsActive returns true if prometheus is activ
func PrometheusIsActive() bool {
	return active
}

type prometheusBackend struct {
	hcl      hcl.Logger
	promSz   map[string]*promSz
	gaugeVec map[string]*prometheus.GaugeVec
	histoVec map[string]*prometheus.HistogramVec
	mu       sync.Mutex
}

type promSz struct {
	Gauges map[string]prometheus.Gauge
}

func (p *prometheusBackend) start() {
	p.hcl.Info("Starting prometheus bridger")
	core.Get().WebServer().Handle("/metrics", promhttp.Handler())
	active = true
}

func (p *prometheusBackend) handleEventBus(e *msg.SzenarioEvtMsg) {
	p.setGauge(e)
	p.setGaugeVec(e)
	// histogram is not useful here
	// p.setHistogramVec(e)
}

const step = "step"

func stepToLabel(s string) string {
	l := s[len(step)+1:]
	l = PrometheusName(l)
	return l
}

func (p *prometheusBackend) setGaugeVec(e *msg.SzenarioEvtMsg) {
	p.mu.Lock()
	defer p.mu.Unlock()
	gv := p.getGaugeVec(e)

	for n, c := range e.Counters {
		if !strings.HasPrefix(n, step) {
			continue
		}
		if f, ok := c.(float64); ok {
			m, err := gv.GetMetricWithLabelValues(stepToLabel(n), PrometheusName(e.Username), PrometheusName(e.Region))
			if err != nil {
				p.hcl.Warnf("GaugeVec: %v", err)
			}
			if m == nil {
				p.hcl.Warnf("GaugeVec %s is nil", n)
				continue
			}
			if errors.Is(e.Err(), szenario.TimeoutError{}) {
				f = downValue
			}
			m.Set(f)
		}
	}
}

func (p *prometheusBackend) getGaugeVec(e *msg.SzenarioEvtMsg) *prometheus.GaugeVec {
	gv := p.gaugeVec[e.Name]
	if gv == nil {
		p.hcl.Infof("Creating GaugeVec %v ", e.Name)
		gv = promauto.NewGaugeVec(prometheus.GaugeOpts{
			//Namespace: PrometheusName(e.Name),
			Name: PrometheusName(e.Name),
			//Subsystem: ,
			// Help: "The total number of processed events",
		},
			[]string{"step", "user", "region"},
		)
		p.gaugeVec[e.Name] = gv
	}
	return gv
}

func (p *prometheusBackend) setGauge(e *msg.SzenarioEvtMsg) {
	p.mu.Lock()
	defer p.mu.Unlock()
	psz := p.getPromSzenario(e)
	for n, c := range e.Counters {
		if f, ok := c.(float64); ok {
			ph := psz.getGauge(e.Name, n)
			if errors.Is(e.Err(), szenario.TimeoutError{}) {
				f = downValue
			}
			ph.Set(f)
		}
	}
}

func (p *prometheusBackend) getPromSzenario(e *msg.SzenarioEvtMsg) *promSz {
	psz := p.promSz[e.Name]
	if psz == nil {
		psz = &promSz{
			Gauges: make(map[string]prometheus.Gauge),
		}
		p.promSz[e.Name] = psz
	}
	return psz
}

// PrometheusName converts a name to a name usable in prometheus
func PrometheusName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return strings.ToLower(s)
}

func (ps *promSz) getGauge(szName, name string) prometheus.Gauge {
	szName = PrometheusName(szName)
	name = PrometheusName(name)
	ph := ps.Gauges[name]
	if ph == nil {
		ph = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: szName,
			Name:      name,
			//Subsystem: ,
			// Help: "The total number of processed events",
		})
		ps.Gauges[name] = ph
	}
	return ph
}

func (p *prometheusBackend) setHistogramVec(e *msg.SzenarioEvtMsg) {
	p.mu.Lock()
	defer p.mu.Unlock()
	hv := p.getHistogramVec(e)

	for n, c := range e.Counters {
		if !strings.HasPrefix(n, step) {
			continue
		}
		if f, ok := c.(float64); ok {
			m, err := hv.GetMetricWithLabelValues(stepToLabel(n), PrometheusName(e.Username), PrometheusName(e.Region))
			if err != nil {
				p.hcl.Warnf("GaugeVec: %v", err)
			}
			if m == nil {
				p.hcl.Warnf("GaugeVec %s is nil", n)
				continue
			}
			if errors.Is(e.Err(), szenario.TimeoutError{}) {
				f = downValue
			}
			m.Observe(f)
		}
	}
}

func (p *prometheusBackend) getHistogramVec(e *msg.SzenarioEvtMsg) *prometheus.HistogramVec {
	name := e.Name + "_histo"
	hv := p.histoVec[name]
	if hv == nil {
		hv = promauto.NewHistogramVec(prometheus.HistogramOpts{
			//Namespace: PrometheusName(e.Name),
			Name: PrometheusName(name),
			//Subsystem: ,
			// Help: "The total number of processed events",
		},
			[]string{"step", "user", "region"},
		)
		p.hcl.Infof("Creating HistogramVec %v ", name)
		p.histoVec[name] = hv
	}
	return hv
}

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
	step      = "step"
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
	//	p.setGauge(e)
	p.setGaugeVec(e)
	// histogram is not useful here
	// p.setHistogramVec(e)
	p.savePasswordInfo(e)
}

func stepToLabel(s string) string {
	l := s[len(step)+1:]
	l = PrometheusName(l)
	return l
}

func (p *prometheusBackend) savePasswordInfo(e *msg.SzenarioEvtMsg) {
	p.saveCounter(e, "logins.passwordage", PrometheusName(e.Username))
	p.saveCounter(e, "logins.failed", PrometheusName(e.Username))
}
func (p *prometheusBackend) saveCounter(e *msg.SzenarioEvtMsg, counterName string, promLabel string) {
	cntr, ok := e.Counters[counterName]
	if !ok {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	gvPwAge := p.getGaugeVecByName(PrometheusName(counterName))

	m, err := gvPwAge.GetMetricWithLabelValues(PrometheusName(e.Name), promLabel, PrometheusName(e.Region))
	if err != nil {
		p.hcl.Warnf("GaugeVec %s: %v", counterName, err)
	}
	if m == nil {
		return
	}
	m.Set(cntr)
}

func (p *prometheusBackend) setGaugeVec(e *msg.SzenarioEvtMsg) {
	p.mu.Lock()
	defer p.mu.Unlock()
	gv := p.getGaugeVec(e)

	for n, f := range e.Counters {
		if !strings.HasPrefix(n, step) {
			continue
		}
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

func (p *prometheusBackend) getGaugeVec(e *msg.SzenarioEvtMsg) *prometheus.GaugeVec {
	return p.getGaugeVecByName(e.Name)
}

func (p *prometheusBackend) getGaugeVecByName(name string) *prometheus.GaugeVec {
	gv := p.gaugeVec[name]
	if gv == nil {
		p.hcl.Infof("Creating GaugeVec %v ", name)
		gv = promauto.NewGaugeVec(prometheus.GaugeOpts{
			//Namespace: PrometheusName(e.Name),
			Name: PrometheusName(name),
			//Subsystem: ,
			// Help: "The total number of processed events",
		},
			[]string{"step", "user", "region"},
		)
		p.gaugeVec[name] = gv
	}
	return gv
}

//nolint:unused
func (p *prometheusBackend) setGauge(e *msg.SzenarioEvtMsg) {
	p.mu.Lock()
	defer p.mu.Unlock()
	psz := p.getPromSzenario(e)
	for n, f := range e.Counters {
		ph := psz.getGauge(e.Name, n)
		if errors.Is(e.Err(), szenario.TimeoutError{}) {
			f = downValue
		}
		ph.Set(f)

	}
}

//nolint:unused
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

//nolint:unused
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

//nolint:unused
func (p *prometheusBackend) setHistogramVec(e *msg.SzenarioEvtMsg) {
	p.mu.Lock()
	defer p.mu.Unlock()
	hv := p.getHistogramVec(e)

	for n, f := range e.Counters {
		if !strings.HasPrefix(n, step) {
			continue
		}
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

//nolint:unused
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

package webstatus

import (
	"fmt"
	"image/color"
	"net/http"
	"strconv"
	"strings"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/vogtp/som/pkg/visualiser/data"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

func (s *WebStatus) handleChart(w http.ResponseWriter, r *http.Request) {
	seq := strings.Split(r.URL.Path, "/")
	szName := seq[len(seq)-1]
	ts, ok := s.data.Timeseries[szName]
	if !ok {
		s.hcl.Warnf("No szenario %s: %v", szName, r.URL.Path)
		http.Error(w, fmt.Sprintf("Szenario %s not found.", szName), http.StatusNotFound)
		return
	}
	//s.doGoChart(w, r, ts)
	s.doGonumPlot(w, r, ts)
}

type dayTickMarks struct{}

func (dayTickMarks) Ticks(min, max float64) []plot.Tick {
	ticks := make([]plot.Tick, 0)
	tmin := time.Unix(int64(min), 0)
	tmax := time.Unix(int64(max), 0)

	day := time.Hour
	tmin = time.Date(tmin.Year(), tmin.Month(), tmin.Day(), 0, 0, 0, 0, time.UTC)

	for ; tmin.Before(tmax); tmin = tmin.Add(day) {
		t := plot.Tick{Value: float64(tmin.Unix())}
		if tmin.Hour() == 0 || tmin.Hour() == 12 {
			t.Label = "x"
		}
		ticks = append(ticks, t)
	}
	// ticks[0].Label = "x"
	// ticks[len(ticks)-1].Label = "x"

	return ticks
}

func (s *WebStatus) doGonumPlot(w http.ResponseWriter, r *http.Request, ts *data.Timeserie) {
	p := plot.New()
	if t := r.URL.Query().Get("title"); len(t) > 0 {
		p.Title.Text = t
	}
	width := 1200
	height := 200
	if w, err := strconv.Atoi(r.URL.Query().Get("width")); err == nil && w > 0 {
		width = w
	}
	if h, err := strconv.Atoi(r.URL.Query().Get("height")); err == nil && h > 0 {
		height = h
	}

	//p.X.Label.Text = "Time"
	p.X.Tick.Marker = plot.TimeTicks{
		Format: "2006-01-02\n15:04",
		Ticker: dayTickMarks{},
	}
	//p.Y.Label.Text = ""
	p.Add(plotter.NewGrid())

	szXYs := func(ts *data.Timeserie) plotter.XYs {

		pts := make(plotter.XYs, ts.Len())
		t, p := ts.GetSlices()
		for i := range t {
			pts[i].X = float64(t[i].Unix())
			pts[i].Y = p[i]
		}
		return pts
	}
	pts := szXYs(ts)

	sc, err := plotter.NewScatter(pts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.hcl.Warnf("error creating plot: %s", err)
		return
	}
	sc.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	sc.GlyphStyle.Radius = vg.Points(1)
	//sc.GlyphStyle.Shape = draw.RingGlyph{}

	p.Add(sc)

	w.Header().Set("Content-Type", "image/svg+xml")
	c := vgsvg.New(vg.Points(float64(width)), vg.Points(float64(height)))
	p.Draw(draw.New(c))
	if _, err := c.WriteTo(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *WebStatus) doGoChart(w http.ResponseWriter, r *http.Request, ts *data.Timeserie) {
	t, p := ts.GetSlices()

	graph := chart.Chart{
		//	Title: szName,
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		Series: []chart.Series{
			chart.TimeSeries{
				//		Style:   chart.Style{StrokeWidth: 0, DotColor: drawing.ColorBlue, DotWidth: 2, DotColorProvider: chart.Shown().DotColorProvider},
				XValues: t,
				YValues: p,
			},
		},
	}

	if t := r.URL.Query().Get("title"); len(t) > 0 {
		graph.Title = t
	}
	if w, err := strconv.Atoi(r.URL.Query().Get("width")); err == nil && w > 0 {
		graph.Width = w
	}
	if h, err := strconv.Atoi(r.URL.Query().Get("height")); err == nil && h > 0 {
		graph.Height = h
	}

	w.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, w)
}

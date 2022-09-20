package webstatus

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

const (
	startParam    = "start"
	endParam      = "end"
	timespanParam = "timespan"

	todayLabel     = "Today"
	yesterdayLabel = "Yesterday"
	last7dLabel    = "Last 7 Days"
	last30dLabel   = "Last 30 Days"
	thisMonthLabel = "This Month"
	lastMonthLabel = "Last Month"
	customLabel    = "Custom Range"
)

type datepicker struct {
	Baseurl  string
	Timespan string
	Start    time.Time
	End      time.Time
}

func (dp *datepicker) init(r *http.Request) {
	dp.Baseurl = core.Get().WebServer().BasePath()
	dp.Timespan = decodedFormGet(r, timespanParam)
	dp.processTimespan(r)
}

func (dp *datepicker) processTimespan(r *http.Request) {
	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()
	switch dp.Timespan {
	case todayLabel:
		t := time.Now()
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end = start.Add(24 * time.Hour).Add(-1 * time.Nanosecond)
	case yesterdayLabel:
		t := time.Now()
		end = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		start = end.Add(-24 * time.Hour)
		end = end.Add(-1 * time.Nanosecond)
	case last7dLabel:
		end = time.Now()
		start = end.Add(-7 * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), start.Location())
	case last30dLabel:
		end = time.Now()
		start = end.Add(-30 * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), start.Location())
	case thisMonthLabel:
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		start = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1).Add(24 * time.Hour).Add(-1 * time.Nanosecond)
	case lastMonthLabel:
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		start = time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1).Add(24 * time.Hour).Add(-1 * time.Nanosecond)
	case customLabel:
		start = parseTime(start, r.Form.Get(startParam))
		end = parseTime(end, r.Form.Get(endParam))
	default:
		dp.Timespan = last7dLabel
		dp.processTimespan(r)
		return
	}
	dp.Start = start
	dp.End = end
	hcl.Infof("Timespan (%s): %s - %s", dp.Timespan, start.Format(cfg.TimeFormatString), end.Format(cfg.TimeFormatString))
}

func decodedFormGet(r *http.Request, key string) string {
	val := r.Form.Get(key)
	dec, err := url.QueryUnescape(val)
	if err == nil {
		return dec
	}
	return val
}

func parseTime(t time.Time, str string) time.Time {
	if len(str) > 0 {
		i, err := strconv.Atoi(str)
		if err == nil {
			return time.Unix(int64(i), 0)
		}
	}
	return t
}
func (datepicker) StartParam() string    { return startParam }
func (datepicker) EndParam() string      { return endParam }
func (datepicker) TimespanParam() string { return timespanParam }

func (datepicker) Today() string     { return todayLabel }
func (datepicker) Yesterday() string { return yesterdayLabel }
func (datepicker) Last7d() string    { return last7dLabel }
func (datepicker) Last30d() string   { return last30dLabel }
func (datepicker) ThisMonth() string { return thisMonthLabel }
func (datepicker) LastMonth() string { return lastMonthLabel }
func (datepicker) Custom() string    { return customLabel }

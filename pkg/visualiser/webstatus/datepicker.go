package webstatus

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/vogtp/som/pkg/core"
)

const (
	startParam    = "start"
	endParam      = "end"
	timespanParam = "timespan"

	todayLabel     = "Today"
	yesterdayLabel = "Yesterday"
	last24hLabel   = "Last 24 Hours"
	last7dLabel    = "Last 7 Days"
	last30dLabel   = "Last 30 Days"
	last90dLabel   = "Last 90 Days"
	last365dLabel  = "Last 365 Days"
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

func (dp datepicker) now() time.Time {
	return time.Now()
}

func (dp *datepicker) processTimespan(r *http.Request) {
	start := dp.now().Add(-30 * 24 * time.Hour)
	end := dp.now()
	switch dp.Timespan {
	case last24hLabel:
		end = dp.now()
		start = end.Add(-24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), start.Location())
	case todayLabel:
		t := dp.now()
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end = start.Add(24 * time.Hour).Add(-1 * time.Nanosecond)
	case yesterdayLabel:
		t := dp.now()
		end = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		start = end.Add(-24 * time.Hour)
		end = end.Add(-1 * time.Nanosecond)
	case last7dLabel:
		end = dp.now()
		start = end.Add(-7 * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), start.Location())
	case last30dLabel:
		end = dp.now()
		start = end.Add(-30 * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), start.Location())
	case last90dLabel:
		end = dp.now()
		start = end.Add(-90 * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), start.Location())
	case last365dLabel:
		end = dp.now()
		start = end.Add(-365 * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), start.Location())
	case thisMonthLabel:
		now := dp.now()
		currentYear, currentMonth, _ := now.Date()
		start = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1).Add(24 * time.Hour).Add(-1 * time.Nanosecond)
	case lastMonthLabel:
		now := dp.now()
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
	//slog.Debug("Timespan (%s): %s - %s", dp.Timespan, start.Format(cfg.TimeFormatString), end.Format(cfg.TimeFormatString))
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

func (datepicker) Last24h() string   { return last24hLabel }
func (datepicker) Today() string     { return todayLabel }
func (datepicker) Yesterday() string { return yesterdayLabel }
func (datepicker) Last7d() string    { return last7dLabel }
func (datepicker) Last30d() string   { return last30dLabel }
func (datepicker) Last90d() string   { return last90dLabel }
func (datepicker) Last365d() string  { return last365dLabel }
func (datepicker) ThisMonth() string { return thisMonthLabel }
func (datepicker) LastMonth() string { return lastMonthLabel }
func (datepicker) Custom() string    { return customLabel }

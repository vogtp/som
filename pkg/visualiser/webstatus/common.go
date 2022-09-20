package webstatus

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

const (
	autoUpdate = "auto_update"
)

type commonData struct {
	Title      string
	Baseurl    string
	Version    string
	Query      string
	AutoReload int
	Start      time.Time
	End        time.Time
}

func common(t string, r *http.Request) *commonData {
	r.ParseForm()
	q := ""
	if len(r.URL.RawQuery) > 0 {
		q = fmt.Sprintf("?%s", r.URL.RawQuery)
	}
	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()
	start = parseTime(start, r.Form.Get("start"))
	end = parseTime(end, r.Form.Get("end"))
	cd := &commonData{
		Title:   t,
		Baseurl: core.Get().WebServer().BasePath(),
		Version: cfg.Version,
		Query:   q,
		Start:   start,
		End:     end,
	}
	if r.Form.Has(autoUpdate) {
		if i, err := strconv.Atoi(r.Form.Get(autoUpdate)); err == nil {
			if i < 60 {
				i = 60
			}
			cd.AutoReload = i
		}
	}
	return cd
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

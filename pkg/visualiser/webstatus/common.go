package webstatus

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/vogtp/som"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

const (
	autoUpdate = "auto_update"
)

type commonData struct {
	Title      string
	TitleImage string
	Baseurl    string
	Version    string
	Query      string
	AutoReload int
	Theme      string
	DatePicker datepicker
}

func common(t string, r *http.Request) *commonData {
	if err := r.ParseForm(); err != nil {
		slog.Warn("Cannot parse form", log.Error, err)
	}
	q := ""
	if len(r.URL.RawQuery) > 0 {
		q = fmt.Sprintf("?%s", r.URL.RawQuery)
	}
	cd := &commonData{
		Title:   t,
		Baseurl: core.Get().WebServer().BasePath(),
		Version: som.Version,
		Query:   q,
		Theme:   "light",
	}
	if theme, err := r.Cookie("theme"); err == nil && theme.Value == "dark" {
		cd.Theme = theme.Value
	}

	cd.DatePicker.init(r)
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

func (commonData) Since(t db.MinMaxTime) time.Duration {
	return time.Since(t.Time()).Truncate(time.Second)
}

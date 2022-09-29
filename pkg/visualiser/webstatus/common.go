package webstatus

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
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
	r.ParseForm()
	q := ""
	if len(r.URL.RawQuery) > 0 {
		q = fmt.Sprintf("?%s", r.URL.RawQuery)
	}
	cd := &commonData{
		Title:   t,
		Baseurl: core.Get().WebServer().BasePath(),
		Version: cfg.Version,
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

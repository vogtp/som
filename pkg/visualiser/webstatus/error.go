package webstatus

import (
	"fmt"
	"net/http"

	"github.com/vogtp/som/pkg/core/log"
)

func (s *WebStatus) Error(w http.ResponseWriter, r *http.Request,
	msg string, err error, status int) {
	w.WriteHeader(status)
	var data = struct {
		*commonData
		Error string
		Msg   string
		URL   string
	}{
		commonData: common("Something went wrong", r),
		Msg:        msg,
		URL:        r.URL.String(),
	}
	if err != nil {
		data.Error = err.Error()
	}
	err2 := templates.ExecuteTemplate(w, "error.gohtml", data)
	if err2 != nil {
		s.log.Error("template error", log.Error, err)
		http.Error(w, fmt.Errorf("%v: %w", err, err2).Error(), http.StatusInternalServerError)
		return
	}
}

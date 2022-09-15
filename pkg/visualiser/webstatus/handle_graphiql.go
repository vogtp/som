package webstatus

import (
	"net/http"

	"github.com/vogtp/som/pkg/core/cfg"
)

func (s *WebStatus) handleGraphiQL(w http.ResponseWriter, r *http.Request) {

	var data = struct {
		*commonData
		TimeFormat string
		MeshInfo   string
	}{
		commonData: common("", r),
		TimeFormat: cfg.TimeFormatString,
	}
	s.render(w, r, "graphiql.gohtml", data)
}

package webstatus

import (
	"net/http"

	"github.com/vogtp/som/pkg/core/cfg"
)

func (s *WebStatus) handleMeshInfo(w http.ResponseWriter, r *http.Request) {

	var data = struct {
		*commonData
		TimeFormat string
	}{
		commonData: common("SOM Mesh", r),
		TimeFormat: cfg.TimeFormatString,
	}
	err := templates.ExecuteTemplate(w, "mesh_info.gohtml", data)
	if err != nil {
		s.hcl.Errorf("Mesh info Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

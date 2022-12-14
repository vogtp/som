package webstatus

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

func (s *WebStatus) handleMeshInfo(w http.ResponseWriter, r *http.Request) {

	var data = struct {
		*commonData
		TimeFormat string
		MeshInfo   string
	}{
		commonData: common("SOM Mesh", r),
		TimeFormat: cfg.TimeFormatString,
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s", core.Get().WebServer().URL(), "bus"))
	if err != nil {
		s.hcl.Errorf("Mesh info request error %v", err)
		s.Error(w, r, "Cannot request mesh info", err, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.hcl.Errorf("Cannot read Mesh info body %v", err)
		s.Error(w, r, "Cannot read mesh info", err, http.StatusInternalServerError)
		return
	}
	data.MeshInfo = string(body)
	s.render(w, r, "mesh_info.gohtml", data)
}

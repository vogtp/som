package webstatus

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
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
		s.log.Error("Mesh info request error", log.Error, err)
		s.Error(w, r, "Cannot request mesh info", err, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Cannot read Mesh info body", log.Error, err)
		s.Error(w, r, "Cannot read mesh info", err, http.StatusInternalServerError)
		return
	}
	data.MeshInfo = string(body)
	s.render(w, r, "mesh_info.gohtml", data)
}

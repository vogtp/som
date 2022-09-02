package webstatus

import (
	"fmt"
	"io/ioutil"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.hcl.Errorf("Cannot read Mesh info body %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.MeshInfo = string(body)

	err = templates.ExecuteTemplate(w, "mesh_info.gohtml", data)
	if err != nil {
		s.hcl.Errorf("Mesh info Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

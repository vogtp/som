package webstatus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vogtp/som/pkg/core/log"
)

const (
	szenarioComponentPath = "/szenario/"
)

func (s *WebStatus) handleSzenarioComponent(w http.ResponseWriter, r *http.Request) {
	szName := ""
	idx := strings.Index(r.URL.Path, szenarioComponentPath)
	if idx > -1 {
		szName = r.URL.Path[idx+len(szenarioComponentPath):]
		for strings.HasSuffix(szName, "/") {
			szName = szName[:len(szName)-1]
		}
		szName = strings.ToLower(szName)
	}
	ext := ""
	if idx := strings.Index(szName, "."); idx > 0 {
		ext = szName[idx+1:]
		szName = szName[:idx]
	}
	accept := r.Header.Get("Accept")
	s.log.Info("Showing szenario compoment", log.Szenario, szName, "extention", ext, "accept_header", accept)
	// FIXME handle mimetypes and extentions
	// browser: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
	// img tag: image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8
	for _, stat := range s.data.Status.Szenarios() {
		if strings.ToLower(stat.Key()) != szName {
			continue
		}

		file := fmt.Sprintf("static/status/%s.png", stat.Level().Img())
		data, err := assetData.ReadFile(file)
		if err != nil {
			s.log.Warn("cannot load status image for", log.Szenario, szName, log.Error, err)
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/png")
		//	w.Header().Add("Content-Type", "image/svg+xml")
		_, err = w.Write(data)
		if err != nil {
			s.log.Warn("Cannot write file", "file", file, log.Error, err)
		}
		return
	}

	s.handleSzenarioComponentInfo(w, r)
}

func (s *WebStatus) handleSzenarioComponentInfo(w http.ResponseWriter, r *http.Request) {
	var data = struct {
		*commonData
		ImgPath   string
		Szenarios []string
	}{
		commonData: common("Szenarios Components", r),
		ImgPath:    szenarioComponentPath,
		Szenarios:  make([]string, len(s.data.Status.Szenarios())),
	}
	for i, stat := range s.data.Status.Szenarios() {
		data.Szenarios[i] = stat.Key()
	}
	s.render(w, r, "szenario_component_info.gohtml", data)
}

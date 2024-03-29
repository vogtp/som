package webstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vogtp/som/pkg/core/log"
)

func (s *WebStatus) render(w http.ResponseWriter, r *http.Request, templateName string, data any) {
	ah := r.Header.Get("Accept")
	s.log.Debug("Render page", "accept_header", ah)

	if strings.Contains(ah, "html") {
		if err := templates.ExecuteTemplate(w, templateName, data); err != nil {
			s.log.Error("cannot render template", "template", templateName, log.Error, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if strings.Contains(ah, "application/json") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.log.Error("cannot encode data to json", "template", templateName, log.Error, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// XML encoder is less flexible than the json encoder -> does not work
	// if strings.Contains(ah, "application/xml") {
	// 	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	// 	if err := xml.NewEncoder(w).Encode(&data); err != nil {
	// 		s.log.Errorf("cannot encode data of %s to xml: %v", templateName, err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	return
	// }

	err := fmt.Errorf("unsupported content-type: %v", ah)
	s.log.Warn("Cannot render", "template", templateName, log.Error, err)
	http.Error(w, err.Error(), http.StatusBadRequest)

}

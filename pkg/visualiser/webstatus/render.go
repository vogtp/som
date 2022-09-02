package webstatus

import "net/http"

func (s *WebStatus) render(w http.ResponseWriter, r *http.Request, templateName string, data any) {

	err := templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		s.hcl.Errorf("cannot render template %s: %v", templateName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

package webstatus

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	// FilesPath is the path for files
	FilesPath = "/files/"
)

func (s *WebStatus) handleFiles(w http.ResponseWriter, r *http.Request) {
	idStr := ""
	idx := strings.Index(r.URL.Path, FilesPath)
	if idx < 1 {
		http.Error(w, "No file ID given", http.StatusBadRequest)
		return
	}
	idStr = strings.ToLower(r.URL.Path[idx+len(FilesPath):])
	if strings.HasSuffix(idStr, "/") {
		idStr = idStr[:len(idStr)-1]
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		s.hcl.Warnf("ID is not a UUID %s: %v", idStr, err)
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}
	s.hcl.Debugf("file %s requested", idStr)

	file, err := s.DB().GetFile(r.Context(), id)
	if err != nil {
		s.hcl.Warnf("No such file %s: %v", idStr, err)
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}
	s.hcl.Debugf("Serving file: %s.%s", file.Name, file.Type.Ext)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", file.Type.MimeType)
	_, err = w.Write(file.Payload)
	if err != nil {
		s.hcl.Warnf("Cannot write file %s: %v", file, err)
	}

}

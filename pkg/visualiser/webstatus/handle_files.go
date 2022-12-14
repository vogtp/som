package webstatus

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/file"
)

const (
	// FilesPath is the path for files
	FilesPath = "/files/"
)

func (s *WebStatus) handleFiles(w http.ResponseWriter, r *http.Request) {
	idStr := ""
	idx := strings.Index(r.URL.Path, FilesPath)
	if idx < 0 {
		http.Error(w, "No file ID given", http.StatusBadRequest)
		return
	}
	idStr = strings.ToLower(r.URL.Path[idx+len(FilesPath):])
	idStr = strings.TrimSuffix(idStr, "/")

	id, err := uuid.Parse(idStr)
	if err != nil {
		s.hcl.Warnf("ID is not a UUID %s: %v", idStr, err)
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}
	s.hcl.Debugf("file %s requested", idStr)

	file, err := s.Ent().File.Query().Where(file.UUIDEQ(id)).First(r.Context())
	if err != nil {
		s.hcl.Warnf("No such file %s: %v", idStr, err)
		http.Error(w, "No such file", http.StatusNotFound)
		return
	}
	s.hcl.Debugf("Serving file: %s.%s", file.Name, file.Ext)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", file.Type)
	_, err = w.Write(file.Payload)
	if err != nil {
		s.hcl.Warnf("Cannot write file %s: %v", file, err)
	}

}

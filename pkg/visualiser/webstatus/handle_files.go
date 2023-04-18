package webstatus

import (
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/log"
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
		s.log.Warn("ID is not a UUID", "id", idStr, log.Error, err)
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}
	s.log.Debug("file requested", "file", idStr)

	file, err := s.Ent().File.Query().Where(file.UUIDEQ(id)).First(r.Context())
	if err != nil {
		s.log.Warn("No such file", "file", idStr, log.Error, err)
		http.Error(w, "No such file", http.StatusNotFound)
		return
	}
	s.log.Debug("Serving file", "file", file.Name, "extention", file.Ext)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", file.Type)
	// _, err = w.Write(file.Payload)
	reader := strings.NewReader(string(file.Payload))
	_, err = io.Copy(w, reader)
	if err != nil {
		s.log.Warn("Cannot write file", "file", file, log.Error, err)
	}

}

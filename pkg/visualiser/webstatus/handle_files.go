package webstatus

import (
	"net/http"
	"strconv"
	"strings"
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

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.hcl.Warnf("ID is not int %s: %v", idStr, err)
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}
	s.hcl.Infof("file %s requested", idStr)

	file, err := s.DB().GetFile(id)
	if err != nil {
		s.hcl.Warnf("No such file %s: %v", idStr, err)
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}
	s.hcl.Infof("Serving file: %s.%s", file.Name, file.Type.Ext)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", file.Type.MimeType)
	_, err = w.Write(file.Payload)
	if err != nil {
		s.hcl.Warnf("Cannot write file %s: %v", file, err)
	}

}

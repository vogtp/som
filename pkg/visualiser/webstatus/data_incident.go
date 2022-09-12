package webstatus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

const (
	incidentfileTimeformat = "20060102150405"
	incidentfileExt        = "json"
	cacheFileName          = "cache.json"
)

func (s *WebStatus) getIncidentRoot() string {
	return fmt.Sprintf("%sincident", s.getDataRoot())
}

func (s *WebStatus) initIncidentTemplates() {
	alerfileStr := fmt.Sprintf(`{{ .Time.Format %q }}-{{.Name}}-{{.ID}}.%s`, incidentfileTimeformat, incidentfileExt)
	s.incidentFileTmpl = template.Must(template.New("incidentFile").Parse(alerfileStr))
	path := fmt.Sprintf(`%s/{{ .Start.Format "2006/01/" }}{{ .Start.Format %q }}-{{.Name}}-{{.IncidentID}}/`, s.getIncidentRoot(), incidentfileTimeformat)
	s.incidentPathTmpl = template.Must(template.New("incidentPath").Parse(path))
	p, err := s.getIncidentPath(msg.NewIncidentMsg(msg.OpenIncident, msg.NewSzenarioEvtMsg("", "", time.Now())))
	if err != nil {
		s.hcl.Errorf("cannot setup incident dumping: %v", err)
		// it is OK to panic here, template.Must panics as well
		panic(err)
	}
	os.Remove(p)
	s.hcl.Infof("saving incidents to %s", s.getIncidentRoot())
}

func (s *WebStatus) saveIncident(a *msg.IncidentMsg) error {
	path, err := s.getIncidentPath(a)
	if err != nil {
		return fmt.Errorf("cannot setup incident path: %v", err)
	}
	if err = s.removeIncidentCache(path); err != nil {
		s.hcl.Warnf("Cannot remove cache file: %v", err)
	}
	fileName, err := s.getIncidentFilename(path, a)
	if err != nil {
		return fmt.Errorf("create incident filename: %v", err)
	}

	d, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("marshalling incident: %v", err)
	}
	s.hcl.Infof("Saving incident of %s to %s (ID: %s)", a.Name, fileName, a.ID)
	return ioutil.WriteFile(fileName, d, 0644)
}

func (s *WebStatus) getIncidentFilename(path string, e *msg.IncidentMsg) (string, error) {
	var f bytes.Buffer
	if err := s.incidentFileTmpl.Execute(&f, e); err != nil {
		return "", fmt.Errorf("cannot execute filename template: %v", err)
	}
	return fmt.Sprintf("%s/%s", path, f.String()), nil
}

func (s *WebStatus) getIncidentPath(e *msg.IncidentMsg) (string, error) {
	var buf bytes.Buffer
	if err := s.incidentPathTmpl.Execute(&buf, e); err != nil {
		return "", fmt.Errorf("cannot execute path template: %v", err)
	}
	p := buf.String()
	if err := core.EnsureOutFolder(p); err != nil {
		return "", err
	}
	return p, nil
}

type incidentInfo struct {
	*msg.IncidentMsg
	Level  status.Level
	Status *statusData
}

type incidentFile struct {
	// Path         string
	// Name         string
	IncidentInfo *incidentInfo
	Error        string
	DetailLink   string
	EvtCnt       uint
}

func (s *WebStatus) getIncidentInfo(file string) (ai *incidentInfo, err error) {

	d, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	im := &msg.IncidentMsg{}
	err = json.Unmarshal(d, im)
	ii := &incidentInfo{
		IncidentMsg: im,
		Level:       status.Level(im.IntLevel),
	}
	if len(ii.ByteState) < 1 {
		s.hcl.Warn("NO BYTE STATE IN incidenInfo")
	}
	return ii, err
}

/*
	func (s *WebStatus) getIncidentFiles(root string, filter string) (fileList []incidentFile, err error) {
		files, err := s.getIncidentDetailFiles(root, filter)
		if err != nil {
			return nil, err
		}
		incidents := make(map[string]incidentFile, 0)
		for _, f := range files {
			i, ok := incidents[f.IncidentInfo.IncidentID]
			if ok {
				if i.IncidentInfo.Start.After(f.IncidentInfo.Start) {
					i.IncidentInfo.Start = f.IncidentInfo.Start
				}
				if i.IncidentInfo.End.Before(f.IncidentInfo.End) {
					i.IncidentInfo.End = f.IncidentInfo.End
				}
				if i.IncidentInfo.End.Before(f.IncidentInfo.Time) {
					i.IncidentInfo.End = time.Time{}
				}
				if i.IncidentInfo.Level < f.IncidentInfo.Level {
					i.IncidentInfo.Level = f.IncidentInfo.Level
				}
				if f.IncidentInfo.Err() != nil {
					i.Error = f.IncidentInfo.Err().Error()
				}
				i.EvtCnt++
				incidents[f.IncidentInfo.IncidentID] = i
				continue
			}
			f.EvtCnt = 1
			incidents[f.IncidentInfo.IncidentID] = f
		}

		ret := make([]incidentFile, 0, len(incidents))
		for _, v := range incidents {
			ret = append(ret, v)
		}
		return ret, nil
	}
*/
func (s *WebStatus) removeIncidentCache(root string) error {
	cacheFile := fmt.Sprintf("%s/%s", root, cacheFileName)
	err := os.Remove(cacheFile)
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err == nil {
		s.hcl.Infof("Removed incident cache file: %s", cacheFile)
	}
	return nil
}

func (s *WebStatus) readIncidentCache(root string) (fileList []incidentFile, err error) {
	cacheFile := fmt.Sprintf("%s/%s", root, cacheFileName)
	cache, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			s.hcl.Infof("Cannot read incident cache %s: %v", cacheFile, err)
		}
		return nil, err
	}
	err = json.Unmarshal(cache, &fileList)
	if err != nil {
		s.hcl.Infof("Cannot unmarshal incident cache %s: %v", cacheFile, err)
		return nil, err
	}
	return fileList, err
}

func (s *WebStatus) writeIncidentCache(root string, fileList []incidentFile) error {
	cacheFile := fmt.Sprintf("%s/%s", root, cacheFileName)
	cache, err := json.Marshal(fileList)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cacheFile, cache, 0644)
}

func filterIncidents(fileList []incidentFile, filter string) []incidentFile {
	if len(filter) < 1 {
		return fileList
	}
	filtered := make([]incidentFile, 0, len(fileList))
	for _, ai := range fileList {
		// if strings.ToLower(ai.Name) != filter &&
		// 	ai.IncidentInfo.IncidentID != filter {
		// 	continue
		// }
		filtered = append(filtered, ai)
	}
	return filtered
}

func (s *WebStatus) getIncidentDetailFiles(a *db.Access, ent *database.Client, root string, filter string) (fileList []incidentFile, err error) {
	// FIXME filter will not work
	// fileList, err = s.readIncidentCache(root)
	// if err == nil {
	// 	return filterIncidents(fileList, filter), nil
	// }

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %v", root, err)
	}
	// baseurl := core.Get().WebServer().BasePath()
	hasFiles := false
	for _, f := range files {
		path := fmt.Sprintf("%s/%s", root, f.Name())
		if f.IsDir() {
			subFiles, err := s.getIncidentDetailFiles(a, ent, path, filter)
			if err != nil {
				return nil, err
			}
			fileList = append(subFiles, fileList...)
			continue
		}
		if f.Name() == cacheFileName {
			continue
		}
		ai, err := s.getIncidentInfo(path)
		if err != nil {
			s.hcl.Warnf("cannot parse incidentinfo from filename %s: %v", f.Name(), err)
			continue
		}
		ctx := context.Background()
		if err := a.SaveIncident(ctx, ai.IncidentMsg); err != nil {
			hcl.Errorf("Save incident: %v", err)
		}

		if err := ent.Incident.Save(ctx, ai.IncidentMsg); err != nil {
			hcl.Warnf("Saving ent incident: %v", err)
		}
		//details := root[len(s.getIncidentRoot())+1:]
		// fileList = append([]incidentFile{
		// 	{
		// 		// Path:         path,
		// 		// Name:         f.Name(),
		// 		IncidentInfo: ai,
		// 		Error:        s.getIncidentError(path),
		// 		DetailLink:   fmt.Sprintf("%s/%s/%s/", baseurl, IncidentDetailPath, ai.IncidentID),
		// 	},
		// }, fileList...)
		hasFiles = true
	}
	if hasFiles {
		if err = s.writeIncidentCache(root, fileList); err != nil {
			s.hcl.Warnf("Could not write incident cache: %v", err)
		}
	}
	return filterIncidents(fileList, filter), nil
}

func (s *WebStatus) getIncidentCount(sz string) int {
	sz = fmt.Sprintf("-%s-", sz)
	return s.countIncidents(s.getIncidentRoot(), sz)
}

func (s *WebStatus) countIncidents(root string, sz string) int {
	c := 0
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return 0
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		c += s.countIncidents(fmt.Sprintf("%s/%s", root, f.Name()), sz)
		if strings.Contains(f.Name(), sz) {
			c++
		}
	}
	return c
}

func (s *WebStatus) getIncidentError(path string) string {
	s.muICache.Lock()
	defer s.muICache.Unlock()
	if e, ok := s.incidentCache[path]; ok {
		return e
	}
	incident, err := s.getIncident(path)
	if err != nil {
		return err.Error()
	}
	if incident.Err() == nil {
		return "OK"
	}
	e := incident.Err().Error()
	s.incidentCache[path] = e
	return e
}

func (s *WebStatus) getIncident(path string) (*incidentInfo, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	im := &msg.IncidentMsg{}
	err = json.Unmarshal(b, im)
	if err != nil {
		return nil, err
	}
	return &incidentInfo{
		IncidentMsg: im,
		Level:       status.Level(im.IntLevel),
	}, nil
}

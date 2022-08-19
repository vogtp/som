package webstatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
	"time"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
)

const (
	alertfileTimeformat = "20060102150405"
	alertfileExt        = "json"
)

func (s *WebStatus) getAlertRoot() string {
	return fmt.Sprintf("%salert", s.getDataRoot())
}

func (s *WebStatus) InitialertTemplates() {
	alerfileStr := fmt.Sprintf(`{{ .Time.Format %q }}-{{.Name}}-{{.ID}}.%s`, alertfileTimeformat, alertfileExt)
	s.alertFileTmpl = template.Must(template.New("alertFile").Parse(alerfileStr))

	path := fmt.Sprintf(`%s/{{ .Format "2006/01/" }}`, s.getAlertRoot())
	s.alertPathTmpl = template.Must(template.New("alertPath").Parse(path))
	_, err := s.getAlertPath(time.Now())
	if err != nil {
		s.hcl.Errorf("cannot setup alert dumping: %v", err)
		// it is OK to panic here, template.Must panics as well
		panic(err)
	}
	s.hcl.Infof("saving alerts to %s", s.getAlertRoot())
}

func (s *WebStatus) saveAlert(a *msg.AlertMsg) error {
	fileName, err := s.getAlertFilename(a)
	if err != nil {
		return fmt.Errorf("create alert filename: %v", err)
	}

	d, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("marshalling alert: %v", err)
	}
	s.hcl.Infof("Saving alert of %s to %s (ID: %s)", a.Name, fileName, a.ID)
	return ioutil.WriteFile(fileName, d, 0644)
}

func (s *WebStatus) getAlertFilename(e *msg.AlertMsg) (string, error) {
	var f bytes.Buffer
	if err := s.alertFileTmpl.Execute(&f, e); err != nil {
		return "", fmt.Errorf("cannot execute filename template: %v", err)
	}
	p, err := s.getAlertPath(e.Time)
	if err != nil {
		return "", fmt.Errorf("cannot setup alert path: %v", err)
	}
	return fmt.Sprintf("%s/%s", p, f.String()), nil
}

func (s *WebStatus) getAlertPath(t time.Time) (string, error) {
	var buf bytes.Buffer
	if err := s.alertPathTmpl.Execute(&buf, t); err != nil {
		return "", fmt.Errorf("cannot execute path template: %v", err)
	}
	p := buf.String()
	if err := core.EnsureOutFolder(p); err != nil {
		return "", err
	}
	return p, nil
}

type alertFile struct {
	Path       string
	Name       string
	AlertInfo  *alertInfo
	Error      string
	DetailLink string
}

type alertInfo struct {
	Time time.Time
	Name string
	ID   string
}

func (s *WebStatus) getAlertInfo(file string) (ai *alertInfo, err error) {
	extLen := len(alertfileExt) + 1 // . is not in the ext
	fileLen := len(file)
	if fileLen < extLen+5 {
		return nil, fmt.Errorf("filename %s is too short", file)
	}
	sub := file[:fileLen-extLen]
	ai = &alertInfo{}
	idx := strings.Index(sub, "-")
	if idx < len(alertfileTimeformat) {
		return nil, fmt.Errorf("time in %s is too shor", file)
	}
	ai.Time, err = time.Parse(alertfileTimeformat, sub[:idx])
	if err != nil {
		return nil, fmt.Errorf("cannot parse time %q in %s: %v", sub[:idx], file, err)
	}
	sub = sub[idx+1:]
	idx = strings.Index(sub, "-")
	if idx < 1 {
		return nil, fmt.Errorf("szenario name in %s is too shor", file)
	}
	ai.Name = sub[:idx]
	ai.ID = sub[idx+1:]
	return ai, nil
}

func (s *WebStatus) getAlertFiles(root string, filter string) (fileList []alertFile, err error) {
	files, err := ioutil.ReadDir(root)
	doFilter := len(filter) > 0
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %v", root, err)
	}
	baseurl := core.Get().WebServer().BasePath()
	for _, f := range files {
		path := fmt.Sprintf("%s/%s", root, f.Name())
		if f.IsDir() {
			subFiles, err := s.getAlertFiles(path, filter)
			if err != nil {
				return nil, err
			}
			fileList = append(subFiles, fileList...)
			continue
		}
		ai, err := s.getAlertInfo(f.Name())
		if err != nil {
			return nil, fmt.Errorf("cannot parse alertinfo from filename %s: %v", f.Name(), err)
		}
		if doFilter &&
			strings.ToLower(ai.Name) != filter &&
			ai.ID != filter {
			continue
		}
		fileList = append([]alertFile{
			{
				Path:       path,
				Name:       f.Name(),
				AlertInfo:  ai,
				Error:      s.getAlertError(path),
				DetailLink: fmt.Sprintf("%s/%s/%s/", baseurl, AlertDetailPath, ai.ID),
			},
		}, fileList...)
	}
	return fileList, nil
}

func (s *WebStatus) getAlertError(path string) string {
	s.muACache.Lock()
	defer s.muACache.Unlock()
	if e, ok := s.alertCache[path]; ok {
		return e
	}
	alert, err := s.getAlert(path)
	if err != nil {
		return err.Error()
	}
	e := alert.Err().Error()
	s.alertCache[path] = e
	return e
}

func (s *WebStatus) getAlert(path string) (*msg.AlertMsg, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	alert := msg.AlertMsg{}
	err = json.Unmarshal(b, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

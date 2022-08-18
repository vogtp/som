package alertmgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vogtp/som/pkg/core/status"
)

const jsonDBFile = "alertmgr.json"

type jsonData struct {
	BasicStates map[string]*basicState
	Status      status.Status
}

func (am *AlertMgr) load() error {
	if err := am.readJSONFile(); err != nil {
		return err
	}
	for n, s := range am.basicStates {
		if s == nil {
			continue
		}
		s.am = am
		s.hcl = am.hcl.Named(n)
	}
	return nil
}

func (am *AlertMgr) save() error {
	return am.wirteJSONFile()
}

func (am *AlertMgr) wirteJSONFile() error {
	am.mu.Lock()
	defer am.mu.Unlock()
	dat := jsonData{
		BasicStates: am.basicStates,
		Status:      am.status,
	}
	b, err := json.Marshal(dat)
	if err != nil {
		return fmt.Errorf("cannot marshal json: %w", err)
	}

	err = ioutil.WriteFile(jsonDBFile, b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file %s: %w", jsonDBFile, err)
	}
	am.hcl.Debugf("Saved %v alertmgr datasets to %s", len(am.basicStates), jsonDBFile)
	return nil
}

func (am *AlertMgr) readJSONFile() error {
	am.mu.Lock()
	defer am.mu.Unlock()
	fi, err := os.Stat(jsonDBFile)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("not found %s: %w", jsonDBFile, err)
	}
	plan, _ := ioutil.ReadFile(fi.Name())

	dat := jsonData{
		BasicStates: am.basicStates,
		Status:      am.status,
	}
	err = json.Unmarshal(plan, &dat)
	if err != nil {
		return fmt.Errorf("error loading json from %v: %w", jsonDBFile, err)
	}
	am.hcl.Infof("Loaded %v alertmgr datasets from %s", len(am.basicStates), fi.Name())
	am.hcl.Debug(am.status.String())
	return nil
}

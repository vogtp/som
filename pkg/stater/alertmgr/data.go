package alertmgr

import (
	"encoding/json"
	"errors"
	"fmt"
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
		s.log = am.log.With("state_name", n)
	}
	status.Cleanup(am.status)
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

	err = os.WriteFile(jsonDBFile, b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file %s: %w", jsonDBFile, err)
	}
	return nil
}

func (am *AlertMgr) readJSONFile() error {
	am.mu.Lock()
	defer am.mu.Unlock()
	fi, err := os.Stat(jsonDBFile)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("not found %s: %w", jsonDBFile, err)
	}
	plan, _ := os.ReadFile(fi.Name())

	dat := jsonData{
		BasicStates: am.basicStates,
		Status:      am.status,
	}
	err = json.Unmarshal(plan, &dat)
	if err != nil {
		return fmt.Errorf("error loading json from %v: %w", jsonDBFile, err)
	}
	am.log.Info("Loaded alertmgr datasets", "count", len(am.basicStates), "file", fi.Name())
	am.log.Debug(am.status.String())
	return nil
}

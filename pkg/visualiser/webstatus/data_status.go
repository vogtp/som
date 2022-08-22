package webstatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/status"
	"github.com/vogtp/som/pkg/visualiser/data"
)

const jsonDBFile = "webstatus.json"

type szenarioData struct {
	hcl           hcl.Logger
	mu            sync.RWMutex `json:"-"`
	Status        status.Status
	Availabilites map[string]status.Availability
	Timeseries    map[string]*data.Timeserie
}

func newSzenarioData(hcl hcl.Logger) *szenarioData {
	sd := &szenarioData{
		hcl:           hcl.Named("data"),
		Status:        status.New(),
		Availabilites: make(map[string]status.Availability),
		Timeseries:    make(map[string]*data.Timeserie),
	}
	return sd
}

func (sz *szenarioData) load() error {
	if err := sz.readJSONFile(); err != nil {
		return err
	}
	status.Cleanup(sz.Status)
	return nil
}

func (sz *szenarioData) save() error {
	if err := sz.wirteJSONFile(); err != nil {
		return err
	}
	return nil
}

func (sz *szenarioData) wirteJSONFile() error {
	sz.mu.Lock()
	defer sz.mu.Unlock()
	b, err := json.MarshalIndent(sz, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal json: %w", err)
	}

	err = ioutil.WriteFile(jsonDBFile, b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file %s: %w", jsonDBFile, err)
	}
	sz.hcl.Debugf("Saved %v szenario datasets to %s", len(sz.Status.Szenarios()), jsonDBFile)
	return nil
}

func (sz *szenarioData) readJSONFile() error {
	sz.mu.Lock()
	defer sz.mu.Unlock()
	fi, err := os.Stat(jsonDBFile)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("not found %s: %w", jsonDBFile, err)
	}
	plan, _ := ioutil.ReadFile(fi.Name())
	err = json.Unmarshal(plan, sz)
	if err != nil {
		return fmt.Errorf("error loading json from %v: %w", jsonDBFile, err)
	}
	sz.hcl.Debugf("Loaded %v szenario datasets from %s", len(sz.Status.Szenarios()), fi.Name())
	return nil
}

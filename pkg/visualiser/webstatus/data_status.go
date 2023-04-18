package webstatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/status"
	"golang.org/x/exp/slog"
)

const jsonDBFile = "webstatus.json"

type szenarioData struct {
	log           *slog.Logger
	mu            sync.RWMutex `json:"-"`
	Status        status.Status
	Availabilites map[string]status.Availability
}

func newSzenarioData(logger *slog.Logger) *szenarioData {
	sd := &szenarioData{
		log:           logger.With(log.Component, "data"),
		Status:        status.New(),
		Availabilites: make(map[string]status.Availability),
	}
	return sd
}

func (sz *szenarioData) load() error {
	if err := sz.readJSONFile(); err != nil {
		return err
	}
	sz.Status.SetConfig(&status.Config{
		UnknownTimeout: viper.GetDuration(cfg.StatusTimeout),
	})
	go func() {
		ticker := time.NewTicker(time.Hour)
		for {
			sz.log.Info("Cleaning status")
			status.Cleanup(sz.Status)
			<-ticker.C
		}
	}()
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

	err = os.WriteFile(jsonDBFile, b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file %s: %w", jsonDBFile, err)
	}
	sz.log.Debug("Saved szenario datasets", "count", len(sz.Status.Szenarios()), "file", jsonDBFile)
	return nil
}

func (sz *szenarioData) readJSONFile() error {
	sz.mu.Lock()
	defer sz.mu.Unlock()
	fi, err := os.Stat(jsonDBFile)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("not found %s: %w", jsonDBFile, err)
	}
	plan, _ := os.ReadFile(fi.Name())
	err = json.Unmarshal(plan, sz)
	if err != nil {
		return fmt.Errorf("error loading json from %v: %w", jsonDBFile, err)
	}
	sz.log.Debug("Loaded szenario datasets", "count", len(sz.Status.Szenarios()), "file", fi.Name())
	return nil
}

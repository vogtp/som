package bridger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"golang.org/x/exp/slog"
)

type ncConfig struct {
	srv, id string
}

// NewNetCrunchBackend creates a backend for NetCrunch Generic Messages
func NewNetCrunchBackend(srv, id string) *ncConfig {
	return &ncConfig{srv: srv, id: id}
}

// RegisterNetCrunchWebMessage registers NetCrunch Messages on the eventbus
func RegisterNetCrunchWebMessage() {
	registerNetCrunchWebMessage("OWA", NewNetCrunchBackend("netcrunch.example.com", "owa@1116"))
}

func registerNetCrunchWebMessage(name string, nc *ncConfig) {
	bus := core.Get().Bus()
	if len(nc.srv) < 1 || len(nc.id) < 1 {
		return
	}
	bus.GetLogger().Info("Got NC backend", "backend", nc)
	ncb := ncBackend{
		log:  bus.GetLogger().With(log.Component, "nc"),
		name: name,
		srv:  nc.srv,
		id:   nc.id,
	}
	bus.Szenario.Handle(ncb.handleEventBus)

}

type ncBackend struct {
	log     *slog.Logger
	name    string
	srv, id string
}

type ncGenericMsg struct {
	Counters map[string]float64 `json:"counters"`
	Statuses map[string]string  `json:"statuses"`
	Retain   int                `json:"retain"`
}

func (nc ncBackend) handleEventBus(e *msg.SzenarioEvtMsg) {
	if nc.name != e.Name {
		// not a msg for us
		return
	}
	hcl := nc.log.With(log.Szenario, e.Name)
	data := ncGenericMsg{
		Counters: e.Counters,
		Statuses: make(map[string]string),
		Retain:   10,
	}

	for k, v := range e.Stati {
		data.Statuses[k] = v
	}

	data.Statuses["error"] = "OK"
	if e.Err() != nil {
		data.Statuses["error"] = fmt.Sprintf("%v", e.Err())
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		hcl.Warn("Cannot encode netcrunch json", log.Error, err)
	}
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/rest/1/sensors/%s/update", nc.srv, nc.id),
		&buf,
	)
	if err != nil {
		hcl.Error("Error creating netcrunch request", log.Error, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		hcl.Error("Error sending to netcrunch", log.Error, err)
		return
	}
	hcl.Info("Sent to netcrunch", "host", nc.srv, "status", resp.Status)
}

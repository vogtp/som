package bridger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
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
	bus.GetLogger().Infof("Got NC backend: %v", nc)
	ncb := ncBackend{
		hcl:  bus.GetLogger().Named("nc"),
		name: name,
		srv:  nc.srv,
		id:   nc.id,
	}
	bus.Szenario.Handle(ncb.handleEventBus)

}

type ncBackend struct {
	hcl     hcl.Logger
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
	hcl := nc.hcl.Named(e.Name)
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
		hcl.Warnf("Cannot encode netcrunch: %v", err)
	}
	hcl.Infof("NC Json: %s to %s", buf.String(), nc.srv)
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/rest/1/sensors/%s/update", nc.srv, nc.id),
		&buf,
	)
	if err != nil {
		hcl.Errorf("Error creating netcrunch request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		hcl.Errorf("Error sending to netcrunch: %v", err)
		return
	}
	hcl.Infof("Sending to %v: %v", nc.srv, resp.Status)
}

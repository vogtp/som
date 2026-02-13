package webstatus

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
)

func (s *WebStatus) handleBusInfo(w http.ResponseWriter, r *http.Request) {

	var data = struct {
		*CommonData
		TimeFormat string
		Connz      Connz
	}{
		CommonData: common("SOM Bus", r),
		TimeFormat: cfg.TimeFormatString,
	}

	resp, err := http.Get("http://localhost:8222/connz?subs=detail&sort=last")
	if err != nil {
		s.log.Error("Bus info request error", log.Error, err)
		s.Error(w, r, "Cannot request bus info", err, http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&data.Connz); err != nil {
		s.log.Error("Cannot parse bus info body", log.Error, err)
		s.Error(w, r, "Cannot parse bus info", err, http.StatusInternalServerError)
		return
	}

	s.render(w, r, "bus_info.gohtml", data)
}

type Connz struct {
	ServerID       string    `json:"server_id"`
	Now            time.Time `json:"now"`
	NumConnections int       `json:"num_connections"`
	Total          int       `json:"total"`
	Offset         int       `json:"offset"`
	Limit          int       `json:"limit"`
	Connections    []Conn    `json:"connections"`
}

type Conn struct {
	Cid                     int       `json:"cid"`
	Kind                    string    `json:"kind"`
	Type                    string    `json:"type"`
	IP                      string    `json:"ip"`
	Port                    int       `json:"port"`
	Start                   time.Time `json:"start"`
	LastActivity            time.Time `json:"last_activity"`
	Rtt                     string    `json:"rtt"`
	Uptime                  string    `json:"uptime"`
	Idle                    string    `json:"idle"`
	PendingBytes            int       `json:"pending_bytes"`
	InMsgs                  int       `json:"in_msgs"`
	OutMsgs                 int       `json:"out_msgs"`
	InBytes                 int       `json:"in_bytes"`
	OutBytes                int       `json:"out_bytes"`
	Subscriptions           int       `json:"subscriptions"`
	SubscriptionsListDetail []struct {
		Subject string `json:"subject"`
		Sid     string `json:"sid"`
		Msgs    int    `json:"msgs"`
		Cid     int    `json:"cid"`
	} `json:"subscriptions_list_detail"`
	Name    string `json:"name"`
	Lang    string `json:"lang"`
	Version string `json:"version"`
}

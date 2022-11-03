package alerter

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/stater/alertmgr"
	"github.com/vogtp/som/pkg/visualiser/webstatus"
)

var (
	//go:embed templates
	assetData embed.FS
	templates = template.Must(template.ParseFS(assetData, "templates/*.gohtml"))
)

func getHTML(e *msg.AlertMsg) (string, error) {
	data := struct {
		Timeformat        string
		AlertDetailURL    string
		IncidentDetailURL string
		Alert             *msg.AlertMsg
		Topology          string
	}{
		Timeformat:        cfg.TimeFormatString,
		AlertDetailURL:    fmt.Sprintf("%s/%s", viper.GetString(cfg.AlertVisualiserURL), webstatus.AlertDetailPath),
		IncidentDetailURL: fmt.Sprintf("%s/%s", viper.GetString(cfg.AlertVisualiserURL), webstatus.IncidentDetailPath),
		Alert:             e,
		Topology:          "none",
	}
	topo, foundTopo := e.Stati[alertmgr.KeyTopology]
	if foundTopo {
		data.Topology = topo
		delete(data.Alert.Stati, alertmgr.KeyTopology)
	}
	var body bytes.Buffer
	err := templates.ExecuteTemplate(&body, "alert.gohtml", &data)
	if foundTopo {
		// we are reusing this alert so readd the topology
		e.Stati[alertmgr.KeyTopology] = topo
	}
	return body.String(), err
}

func getSubject(e *msg.AlertMsg, r *Rule, d *Destination) string {
	subj := fmt.Sprintf("%s %s - ", getCfgString(cfgAlertSubject, r, d), e.Level)
	if !strings.Contains(strings.ToLower(e.Err().Error()), strings.ToLower(e.Name)) {
		subj = fmt.Sprintf("%v %s:", subj, e.Name)
	}
	return fmt.Sprintf("%v %s", subj, e.Err())
}

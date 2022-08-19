package alerter

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/mime"
	"github.com/vogtp/som/pkg/core/msg"
)

type altr interface {
	sendAlert(e *msg.AlertMsg) error
}
type alertTestData struct {
	name    string
	alerter altr
}

func Test_sendAlert(t *testing.T) {
	tests := make([]alertTestData, 0)
	teamsWebHook := viper.GetString(cfg.AlertTeamsWebhook)
	if len(teamsWebHook) > 0 {
		tests = append(tests, alertTestData{
			name: "teams",
			alerter: &Teams{
				hcl:        hcl.New(),
				webhookURL: teamsWebHook,
			},
		})
	}
	smtpHost := viper.GetString(cfg.AlertMailSMTPHost)
	if len(smtpHost) > 0 {
		tests = append(tests, alertTestData{
			name: "mail",
			alerter: &Mail{
				hcl:      hcl.New(),
				smtpHost: smtpHost,
				smtpPort: viper.GetInt(cfg.AlertMailSMTPPort),
				to:       viper.GetStringSlice(cfg.AlertMailTo),
				from:     viper.GetString(cfg.AlertMailFrom),
			},
		})
	}

	evt := msg.NewSzenarioEvtMsg("unit test", "test user", time.Now())
	evt.AddErr(errors.New("unit test alert"))
	evt.SetCounter("test counter", 42)
	evt.SetCounter("example counter", 2.7165432156789)
	evt.SetStatus("test status", "unit test")
	img, err := ioutil.ReadFile("./testData/screenshot.png")
	if err != nil {
		t.Fatalf("cannot read test image: %v", err)
	}
	evt.AddFile(msg.NewFileMsgItem("unit test screenshot", mime.Png, img))
	al := msg.NewAlert(evt)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt.SetStatus("testcase", tt.name)
			err := tt.alerter.sendAlert(al)
			if err != nil {
				t.Errorf("cannot send to %s: %v", tt.name, err)
			}
		})
	}
}

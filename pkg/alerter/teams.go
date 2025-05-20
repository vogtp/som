package alerter

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"sync"

	"log/slog"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/nfnt/resize"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
)

const (
	targetURLDesc = "SOM Homepage"
)

// Teams  alerter
type Teams struct {
	log *slog.Logger
	mu  sync.Mutex
}

// NewTeams registers a Teams alerter on the event bus
func NewTeams() (Engine, error) {
	bus := core.Get().Bus()
	log := bus.GetLogger().With(log.Component, "teams")
	return &Teams{
		log: log,
	}, nil
}

// Kind returns what kind of alerter engine it is
func (teams *Teams) Kind() string { return "teams" }

// Send the alert
func (teams *Teams) Send(e *msg.AlertMsg, r *Rule, d *Destination) error {
	teams.mu.Lock()
	defer teams.mu.Unlock()
	teams.log.Info("Sending teams alert", "alert", e.Name, "message", e.Err(), "rule", r.name, "destination", d.name)
	err := teams.sendAlert(e, r, d)
	if err != nil {
		return fmt.Errorf("cannot send to teams: %v", err)
	}
	return nil
}

func (teams *Teams) checkConfig(a *Alerter) (ret error) {
	for _, r := range a.rules {
		r := r
		for _, d := range r.destinations {
			d := d
			if d.kind != teams.Kind() {
				continue
			}
			url := d.cfg.GetString(cfgAlertDestTeamsWebhook)

			if err := goteamsnotify.NewClient().ValidateWebhook(url); err != nil {
				teams.log.Warn("teams webhook URL not valid", "destination", d.name, "rule", r.name, "webhook", url, log.Error, err)
				ret = err
			}
			if len(getCfgString(cfgAlertSubject, &r, &d)) < 1 {
				teams.log.Warn("no subject for teams", "destination", d.name, "rule", r.name)
			}
		}
	}
	return ret
}

func (teams *Teams) sendAlert(e *msg.AlertMsg, r *Rule, d *Destination) error {
	mstClient := goteamsnotify.NewClient()
	webhookURL := d.cfg.GetString(cfgAlertDestTeamsWebhook)
	if err := goteamsnotify.NewClient().ValidateWebhook(webhookURL); err != nil {
		return fmt.Errorf("not sending teams message webhook URL %s not valid: %v", webhookURL, err)
	}
	img := ""
	//fullImage := fmt.Sprintf("<img src='%s/%s?file=%%s.%%s' />", viper.GetString(cfg.AlertDetailURL), e.ID)
	// for _, f := range e.Files {
	// 	name := fmt.Sprintf("%s.%s", f.Name, f.Type.Ext)
	// 	teams.log.Debugf("Adding attachment %s", name)
	// 	if f.Type == mime.Png {
	// 		imgTag, err := teams.getImage(f.Payload)
	// 		if err != nil {
	// 			teams.log.Warnf("cannot add image %s: %v", name, err)
	// 			continue
	// 		}
	// 		//teams.log.Infof("%s",fmt.Sprintf(fullImage, f.Name, f.Type.Ext))
	// 		//img = fmt.Sprintf("%s\n%s<br />\n", img, fmt.Sprintf(fullImage, f.Name, f.Type.Ext))
	// 		img = fmt.Sprintf("%s\n%s<br />\n", img, imgTag)
	// 		break
	// 	}
	// }
	text, err := getHTML(e)
	if err != nil {
		teams.log.Error("index Template error", log.Error, err, "destination", d.name, "rule", r.name)
	}
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = getSubject(e, r, d)
	msgCard.Text = fmt.Sprintf("<p>\n%s\n</p>\n%s", text, img)

	pa, err := goteamsnotify.NewMessageCardPotentialAction(
		goteamsnotify.PotentialActionOpenURIType,
		targetURLDesc,
	)

	if err != nil {
		return fmt.Errorf("error creating new teams action: %w", err)
	}
	pa.MessageCardPotentialActionOpenURI.Targets =
		[]goteamsnotify.MessageCardPotentialActionOpenURITarget{
			{
				OS:  "default",
				URI: viper.GetString(cfg.AlertVisualiserURL),
			},
		}
	pa.MessageCardPotentialActionOpenURI.Targets = nil
	if err := msgCard.AddPotentialAction(pa); err != nil {
		return fmt.Errorf("error creating new teams message card: %w", err)
	}
	if err := msgCard.Validate(); err != nil {
		teams.log.Error("msg card is not valid", log.Error, err, "destination", d.name, "rule", r.name)
		return fmt.Errorf("msg card is not valid: %w", err)
	}
	return mstClient.SendWithRetry(context.TODO(), webhookURL, msgCard, 3, 5)
	//return mstClient.Send(webhookURL, msgCard)
}

//nolint:unused
func (teams *Teams) getImage(img []byte) (string, error) {
	image, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return "", fmt.Errorf("cannot decode image: %w", err)
	}

	newImage := resize.Resize(100, 0, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	err = enc.Encode(buf, newImage)
	if err != nil {
		return "", fmt.Errorf("cannot encode image to png: %w", err)
	}

	imgb64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	teams.log.Debug("image len", "len", len(imgb64))

	return fmt.Sprintf("<img src='data:image/png;base64, %s' />", imgb64), nil
}

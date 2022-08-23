package alerter

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"sync"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/nfnt/resize"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/mime"
	"github.com/vogtp/som/pkg/core/msg"
)

const (
	targetURLDesc = "SOM Homepage"
)

// Teams  alerter
type Teams struct {
	hcl hcl.Logger
	mu  sync.Mutex
}

// NewTeams registers a Teams alerter on the event bus
func NewTeams() (Engine, error) {
	bus := core.Get().Bus()
	hcl := bus.GetLogger().Named("teams")
	return &Teams{
		hcl: hcl,
	}, nil
}

// Kind returns what kind of alerter engine it is
func (teams *Teams) Kind() string { return "teams" }

// Send the alert
func (teams *Teams) Send(e *msg.AlertMsg, d *Destination) error {
	teams.mu.Lock()
	defer teams.mu.Unlock()
	teams.hcl.Debug("got event %v: %v", e.Name, e.Err())
	teams.hcl.Infof("Sending teams alert %s: %v", e.Name, e.Err())
	err := teams.sendAlert(e, d)
	if err != nil {
		return fmt.Errorf("cannot send to teams: %v", err)
	}
	return nil
}

func (teams *Teams) checkDestinationWebhooks() (ret error) {
	for _, d := range destinations {
		if d.Kind != teams.Kind() {
			continue
		}
		url := d.Cfg.GetString(cfgAlertDestTeamsWebhook)
		if ok, err := goteamsnotify.IsValidWebhookURL(url); !ok || err != nil {
			teams.hcl.Warnf("%s: teams webhook URL %q not valid: %v", d.Name, url, err)
			if err != nil {
				ret = err
			}
		}
	}
	return ret
}

func (teams *Teams) sendAlert(e *msg.AlertMsg, d *Destination) error {
	webhookURL := d.Cfg.GetString(cfgAlertDestTeamsWebhook)
	if ok, err := goteamsnotify.IsValidWebhookURL(webhookURL); !ok || err != nil {
		return fmt.Errorf("not sending teams message webhook URL %s not valid: %v", webhookURL, err)
	}
	img := ""
	//fullImage := fmt.Sprintf("<img src='%s/%s?file=%%s.%%s' />", viper.GetString(cfg.AlertDetailURL), e.ID)
	for _, f := range e.Files {
		name := fmt.Sprintf("%s.%s", f.Name, f.Type.Ext)
		teams.hcl.Debugf("Adding attachment %s", name)
		if f.Type == mime.Png {
			imgTag, err := teams.getImage(f.Payload)
			if err != nil {
				teams.hcl.Warnf("cannot add image %s: %v", name, err)
				continue
			}
			//teams.hcl.Infof("%s",fmt.Sprintf(fullImage, f.Name, f.Type.Ext))
			//img = fmt.Sprintf("%s\n%s<br />\n", img, fmt.Sprintf(fullImage, f.Name, f.Type.Ext))
			img = fmt.Sprintf("%s\n%s<br />\n", img, imgTag)
			break
		}
	}
	text, err := getHTML(e)
	if err != nil {
		teams.hcl.Errorf("index Template error %v", err)
	}
	mstClient := goteamsnotify.NewClient()
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = getSubject(e)
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

	return mstClient.SendWithRetry(context.TODO(), webhookURL, msgCard, 3, 5)
}

func (teams *Teams) getImage(img []byte) (string, error) {
	image, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return "", fmt.Errorf("cannot decode image: %w", err)
	}

	newImage := resize.Resize(150, 0, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	err = enc.Encode(buf, newImage)
	if err != nil {
		return "", fmt.Errorf("cannot encode image to png: %w", err)
	}

	imgb64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	teams.hcl.Debugf("image len: %v", len(imgb64))

	return fmt.Sprintf("<img src='data:image/png;base64, %s' />", imgb64), nil
}

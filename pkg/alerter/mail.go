package alerter

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"gopkg.in/gomail.v2"
)

// Mail is a mail alerter
type Mail struct {
	hcl      hcl.Logger
	smtpHost string
	smtpPort int
	to       []string
	from     string
}

// NewMailer registers a mail alerter on the event bus
func NewMailer(to ...string) {
	bus := core.Get().Bus()
	hcl := bus.GetLogger().Named("mail")
	mailHost := viper.GetString(cfg.AlertMailSMTPHost)
	if len(mailHost) < 1 {
		hcl.Errorf("Not starting mail alerter: no mail host given")
		return
	}
	if len(viper.GetStringSlice(cfg.AlertMailTo)) < 1 {
		viper.SetDefault(cfg.AlertMailTo, to)
	}
	to = append(to, viper.GetStringSlice(cfg.AlertMailTo)...)
	if len(to) < 1 {
		hcl.Errorf("Not starting mail alerter: no mail receipients given")
		return
	}
	d := Mail{
		hcl:      hcl,
		smtpHost: mailHost,
		smtpPort: viper.GetInt(cfg.AlertMailSMTPPort),
		to:       to,
		from:     viper.GetString(cfg.AlertMailFrom),
	}
	bus.Alert.Handle(d.handle)
}

func (alt *Mail) handle(e *msg.AlertMsg) {
	alt.hcl.Debug("got event %v: %v", e.Name, e.Err())
	err := alt.sendAlert(e)
	if err != nil {
		alt.hcl.Errorf("cannot send mail: %v", err)
	}
}

func (alt *Mail) attachFile(f *msg.FileMsgItem) gomail.FileSetting {
	return gomail.SetCopyFunc(func(w io.Writer) error {
		s, err := w.Write(f.Payload)
		if err != nil {
			alt.hcl.Warnf("cannot attach file %q: %v", f.Name, err)
		}
		if s != f.Size {
			alt.hcl.Warnf("not written enough %v should be %v", s, f.Size)
		}
		return err
	})
}

func (alt *Mail) sendAlert(e *msg.AlertMsg) error {
	if len(alt.to) < 1 {
		alt.hcl.Debugf("No mail-to not sending %s: %v", e.Name, e.Err())
		return nil
	}
	subj := getSubject(e)
	m := gomail.NewMessage()
	m.SetHeader("From", alt.from)
	m.SetHeader("To", alt.to...)
	m.SetHeader("Subject", subj)
	img := ""
	for _, f := range e.Files {
		name := fmt.Sprintf("%s.%s", f.Name, f.Type.Ext)
		alt.hcl.Debugf("Adding attachment %s", name)
		header := make(map[string][]string)
		header["Content-Type"] = []string{f.Type.MimeType}
		if strings.HasPrefix(f.Type.MimeType, "image/") {
			m.Embed(name, alt.attachFile(&f), gomail.SetHeader(header))
			m.Attach(name, alt.attachFile(&f), gomail.SetHeader(header))
			img = fmt.Sprintf(`%s<br /> <img src="cid:%s" alt="My image" />`, img, name)
		} else {
			m.Attach(name, alt.attachFile(&f), gomail.SetHeader(header))
		}
	}

	body, err := getHTML(e)
	if err != nil {
		return fmt.Errorf("index Template error %v", err)
	}

	bd := fmt.Sprintf("<html><body>%s%s<br /><small>SOM Version: %s</small></body></html>", body, img, cfg.Version)
	m.SetBody("text/html", bd)

	d := gomail.Dialer{Host: alt.smtpHost, Port: alt.smtpPort}
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("cannot send mail: %w", err)
	}
	alt.hcl.Infof("Sent email %q to %v", subj, alt.to)
	return nil
}

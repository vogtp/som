package alerter

import (
	"fmt"
	"io"
	"strings"
	"sync"

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
	mu       sync.Mutex
	smtpHost string
	smtpPort int
	//	to       []string
	from string
}

// NewMailer registers a mail alerter on the event bus
func NewMailer() (Engine, error) {
	bus := core.Get().Bus()
	hcl := bus.GetLogger().Named("mail")
	mailHost := viper.GetString(cfg.AlertMailSMTPHost)
	if len(mailHost) < 1 {
		return nil, fmt.Errorf("not creating mail alerter: no mail host given")
	}
	return &Mail{
		hcl:      hcl,
		smtpHost: mailHost,
		smtpPort: viper.GetInt(cfg.AlertMailSMTPPort),
		from:     viper.GetString(cfg.AlertMailFrom),
	}, nil
}

// Kind returns what kind of alerter engine it is
func (alt *Mail) Kind() string { return "mail" }

// Send the alert
func (alt *Mail) Send(e *msg.AlertMsg, r *Rule, d *Destination) error {
	alt.mu.Lock()
	defer alt.mu.Unlock()
	alt.hcl.Debug("got event %v: %v", e.Name, e.Err())

	if err := alt.sendAlert(e, r, d); err != nil {
		return fmt.Errorf("cannot send mail: %v", err)
	}
	return nil
}

func (alt *Mail) checkConfig(a *Alerter) (ret error) {
	if !strings.Contains(alt.from, "@") {
		ret = fmt.Errorf("mail from %q is not valid", alt.from)
		alt.hcl.Warn(ret.Error())
	}
	for _, d := range a.dsts {
		if d.Kind != alt.Kind() {
			continue
		}
		to := d.Cfg.GetStringSlice(cfgAlertDestMailTo)
		if len(to) < 1 {
			ret = fmt.Errorf("mail dest %q: no receipients %q", d.Name, to)
			alt.hcl.Warn(ret.Error())
		}
		for _, t := range to {
			if !strings.Contains(t, "@") {
				ret = fmt.Errorf("mail dest %q: invaluid to %q", d.Name, t)
				alt.hcl.Warn(ret.Error())
			}
		}
	}
	return ret
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

func (alt *Mail) sendAlert(e *msg.AlertMsg, r *Rule, d *Destination) error {
	to := d.Cfg.GetStringSlice(cfgAlertDestMailTo)
	if len(to) < 1 {
		alt.hcl.Debugf("No mail-to not sending %s: %v", e.Name, e.Err())
		return nil
	}
	subj := getSubject(e, r, d)
	m := gomail.NewMessage()
	m.SetHeader("From", alt.from)
	m.SetHeader("To", to...)
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

	mailer := gomail.Dialer{Host: alt.smtpHost, Port: alt.smtpPort}
	if err := mailer.DialAndSend(m); err != nil {
		return fmt.Errorf("cannot send mail: %w", err)
	}
	alt.hcl.Infof("Sent email %q to %v", subj, to)
	return nil
}

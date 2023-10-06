package alerter

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"log/slog"

	"github.com/spf13/viper"
	"github.com/vogtp/som"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"gopkg.in/gomail.v2"
)

// Mail is a mail alerter
type Mail struct {
	log      *slog.Logger
	mu       sync.Mutex
	smtpHost string
	smtpPort int
	from     string
}

// NewMailer registers a mail alerter on the event bus
func NewMailer() (Engine, error) {
	bus := core.Get().Bus()
	log := bus.GetLogger().With("alerter", "mail")
	mailHost := viper.GetString(cfg.AlertMailSMTPHost)
	if len(mailHost) < 1 {
		return nil, fmt.Errorf("not creating mail alerter: no mail host given")
	}
	return &Mail{
		log:      log,
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
	alt.log.Debug("got event %v: %v", e.Name, e.Err())

	if err := alt.sendAlert(e, r, d); err != nil {
		return fmt.Errorf("cannot send mail: %v", err)
	}
	return nil
}

func (alt *Mail) checkConfig(a *Alerter) (ret error) {
	if !strings.Contains(alt.from, "@") {
		ret = fmt.Errorf("mail from %q is not valid", alt.from)
		alt.log.Warn(ret.Error())
	}
	for _, r := range a.rules {
		r := r
		for _, d := range r.destinations {
			d := d
			if d.kind != alt.Kind() {
				continue
			}
			if len(getCfgString(cfgAlertSubject, &r, &d)) < 1 {
				alt.log.Warn("mail has no subject", "rule", r.name, "destination", d.name)
			}
			to := d.cfg.GetStringSlice(cfgAlertDestMailTo)
			if len(to) < 1 {
				ret = fmt.Errorf("mail dest %q: no receipients %q", d.name, to)
				alt.log.Warn(ret.Error())
			}
			for _, t := range to {
				if !strings.Contains(t, "@") {
					ret = fmt.Errorf("mail dest %q: invaluid to %q", d.name, t)
					alt.log.Warn(ret.Error())
				}
			}
		}
	}
	return ret
}

func (alt *Mail) attachFile(f *msg.FileMsgItem) gomail.FileSetting {
	return gomail.SetCopyFunc(func(w io.Writer) error {
		s, err := w.Write(f.Payload)
		if err != nil {
			alt.log.Warn("cannot attach file", "file", f.Name, log.Error, err)
		}
		if s != f.Size {
			alt.log.Warn("not written enough", "file", f.Name, "bytes_wirtten", s, "bytes_total", f.Size)
		}
		return err
	})
}

func (alt *Mail) sendAlert(e *msg.AlertMsg, r *Rule, d *Destination) error {
	to := d.cfg.GetStringSlice(cfgAlertDestMailTo)
	if len(to) < 1 {
		alt.log.Debug("No mail-to not sending", "alert", e.Name, "message", e.Err(), "rule", r.name, "destination", d.name)
		return nil
	}
	subj := getSubject(e, r, d)
	m := gomail.NewMessage()
	m.SetHeader("From", alt.from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subj)
	img := ""
	for _, f := range e.Files {
		f := f
		name := fmt.Sprintf("%s.%s", f.Name, f.Type.Ext)
		alt.log.Debug("Adding attachment", "attachment", name)
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

	bd := fmt.Sprintf("<html><body>%s%s<br /><small>SOM Version: %s</small></body></html>", body, img, som.Version)
	m.SetBody("text/html", bd)

	mailer := gomail.Dialer{Host: alt.smtpHost, Port: alt.smtpPort}
	if err := mailer.DialAndSend(m); err != nil {
		return fmt.Errorf("cannot send mail: %w", err)
	}
	alt.log.Info("Sent email", "subject", subj, "destination", to, "rule", r.name, "destination", d.name)
	return nil
}

package szenarios

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

func AzOIDCLogin(engine szenario.Engine, s szenario.Szenario) (err error) {
	user := s.User()
	err = engine.StepTimeout("Enter EMail", 20*time.Second,
		chromedp.WaitReady("[name='loginfmt']", chromedp.ByQuery),
		chromedp.SendKeys("[name='loginfmt']", fmt.Sprintf("%s\r", user.Email()), chromedp.ByQuery),
		//chromedp.Click("[id='idSIButton9']", chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Errorf("cannot enter email: %w", err)
	}
	err = engine.StepTimeout("Enter Password", 20*time.Second,
		chromedp.WaitReady("[name='Password']", chromedp.ByQuery),
		chromedp.SendKeys("[name='Password']", fmt.Sprintf("%s\r", user.Password()), chromedp.ByQuery),
		// chromedp.Click("[id='submitButton']", chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Errorf("cannot enter Password: %w", err)
	}
	return nil
}

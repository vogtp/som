package szenarios

import (
	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

func AdfsLogin(engine szenario.Engine, s szenario.Szenario) (err error) {

	engine.Step("Login",
		chromedp.WaitVisible(`#userNameInput`, chromedp.ByID),
		chromedp.SendKeys(`#userNameInput`, s.User().Name()+"\r", chromedp.ByID),
		chromedp.WaitVisible(`#passwordInput`, chromedp.ByID),
		chromedp.WaitEnabled(`#passwordInput`, chromedp.ByID),
		chromedp.SendKeys(`#passwordInput`, s.User().Password()+"\r", chromedp.ByID),
	)

	return nil
}

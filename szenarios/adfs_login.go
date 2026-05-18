package szenarios

import (
	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

func AdfsLogin(engine szenario.Engine, s szenario.Szenario) (err error) {

	userInputID := `#userNameInput`
	pwInputID := `#passwordInput`
	engine.Step("Login",
		chromedp.WaitVisible(userInputID, chromedp.ByID),
		chromedp.SendKeys(userInputID, s.User().Name()+"\r", chromedp.ByID),
		chromedp.WaitVisible(pwInputID, chromedp.ByID),
		chromedp.WaitEnabled(pwInputID, chromedp.ByID),
		chromedp.SendKeys(pwInputID, s.User().Password()+"\r", chromedp.ByID),
	)

	trustBuID := `#idSIButton9`
	trust := "trust"
	notrust := "notrust"
	opt := engine.Either("Is trust displayed",
		szenario.EitherOption{ID: trust, Action: chromedp.WaitVisible(trustBuID, chromedp.ByID)},
		szenario.EitherOption{ID: notrust, Action: chromedp.WaitNotVisible(trustBuID, chromedp.ByID)},
	)

	sel := <-opt
	if sel == trust {
		engine.Step("Trust ourself",
			chromedp.Click(trustBuID, chromedp.ByID),
			// chromedp.WaitNotVisible(trustBuID, chromedp.ByID),
		)
	}

	return nil
}

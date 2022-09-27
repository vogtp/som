package szenarios

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// OwaSzenario load Outlook Web Acces (on prem) and checks that the inbox is visible
type OwaSzenario struct {
	*szenario.Base
	OwaURL string
}

// Execute the szenario
func (s *OwaSzenario) Execute(engine szenario.Engine) (err error) {
	engine.Step("Loading",
		chromedp.Navigate(s.OwaURL),
		chromedp.WaitVisible(`#userNameInput`, chromedp.ByID),
	)
	engine.Step("Login",
		chromedp.WaitVisible(`#userNameInput`, chromedp.ByID),
		chromedp.SendKeys(`#userNameInput`, s.User().Name()+"\r", chromedp.ByID),
		chromedp.WaitReady(`#passwordInput`, chromedp.ByID),
		chromedp.SendKeys(`#passwordInput`, s.User().Password()+"\r", chromedp.ByID),
	)
	defer func() {
		engine.Step("Logout",
			chromedp.Navigate(s.OwaURL+"/owa/logoff.owa"),
			chromedp.WaitVisible(`#openingMessage`, chromedp.ByID),
		)
	}()

	loadedID := `#O365_MainLink_NavMenu,#EndOfLifeMessageDiv`
	
	engine.Step("Wait for loading", chromedp.WaitReady(loadedID+",.errorHeader", chromedp.ByID))

	if !engine.IsPresent(loadedID, chromedp.ByID) {
		engine.AddErr(fmt.Errorf("Error page loaded"))
		return
	}

	engine.Step("check loaded", engine.Body(engine.Contains("Sent Items"), engine.Contains("Inbox"), engine.Bigger(1000)))

	return nil
}

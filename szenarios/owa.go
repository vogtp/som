package szenarios

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// OwaSzenario load Outlook Web Acces (on prem) and checks that the inbox is visible
type OwaSzenario struct {
	*szenario.Base
	OwaURL       string
	LoginTimeout time.Duration
}

// Execute the szenario
func (s *OwaSzenario) Execute(engine szenario.Engine) (err error) {

	loggedIn := false
	for !loggedIn {
		err := s.login(engine)
		if err == nil {
			loggedIn = true
		}
		if s.User().NextPassword() == "" {
			return err
		}
	}

	//engine.WaitForEver()

	if engine.IsHeadless() {
		defer func() {
			engine.Step("Logout",
				chromedp.Navigate(s.OwaURL+"/owa/logoff.owa"),
			//	chromedp.WaitVisible(`#openingMessage`, chromedp.ByID),
			)
		}()
	}

	engine.Step("check loaded", engine.Body(engine.Contains("Sent Items"), engine.Contains("Inbox"), engine.Bigger(100)))
	return nil
}

func (s *OwaSzenario) login(engine szenario.Engine) (err error) {
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
	loadedID := `#O365_MainLink_NavMenu,#EndOfLifeMessageDiv`
	if s.LoginTimeout == 0 {
		s.LoginTimeout = 2 * time.Second
	}
	if err := engine.StepTimeout("Check for loading", s.LoginTimeout,
		chromedp.WaitReady(loadedID+",.errorHeader", chromedp.ByID),
	); err != nil {
		//set language and timezone manually
		return fmt.Errorf("Set language and timezone manually and this will be gone: %v", err)
		// //*[@id="selTz"]/option[49]
		// engine.Step("Accept first login",
		// 	chromedp.Click(`*[@id="selTz"]/option[49]`, chromedp.BySearch),
		// 	//	chromedp.Click(`//select[@id="selTz"]/option[@value="UTC"]`, chromedp.BySearch),
		// 	//	chromedp.Click(`document.querySelector("select#selTz > option[value="UTC"]")`, chromedp.ByQuery),
		// 	chromedp.Click(`document.querySelector("#lgnDiv > div.hidden-submit > input[type=submit]")`, chromedp.ByJSPath),
		// 	chromedp.WaitNotVisible("#selTz", chromedp.ByID),
		// )
		// engine.Step("Wait for loading", chromedp.WaitReady(loadedID+",.errorHeader", chromedp.ByID))
	}

	if !engine.IsPresent(loadedID, chromedp.ByID) {
		err = fmt.Errorf("Error page loaded")
		engine.AddErr(err)
		return err
	}
	return err
}

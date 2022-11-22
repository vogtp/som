package szenarios

import (
	"errors"

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
	)
	loadedID := `#O365_MainLink_NavMenu,#EndOfLifeMessageDiv,#lnkHdrcheckmessages`
	errorClass := `#error`
	submitBu := `#nextButton`

	for i := 0; i < s.MaxLoginTry(); i++ {
		engine.Step("Login",
			chromedp.WaitReady(`#passwordInput`, chromedp.ByID),
			chromedp.SendKeys(`#passwordInput`, s.User().Password()+"\r", chromedp.ByID),
		)
		option := <-engine.Either("Login Check",
			szenario.EitherOption{ID: loadedID, Action: chromedp.WaitVisible(loadedID, chromedp.ByID)},
			szenario.EitherOption{ID: errorClass, Action: chromedp.WaitVisible(errorClass, chromedp.ByID)},
		)
		if option == loadedID {
			return nil
		}
		engine.Step("Login",
			chromedp.WaitReady(submitBu, chromedp.ByID),
			chromedp.Click(submitBu, chromedp.ByID),
		)
		pw := s.User().NextPassword()
		if pw == "" {
			break
		}
	}
	return errors.New("Login failed")
}

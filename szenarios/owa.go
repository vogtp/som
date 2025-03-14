package szenarios

import (
	"errors"
	"time"

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

	if err := s.login(engine); err != nil {
		return err
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
	loadedID := `#O365_MainLink_NavMenu,#EndOfLifeMessageDiv,#lnkHdrcheckmessages,#lnkNavMail`
	errorClass := `#error`

	for range s.GetMaxLoginTry() {
		engine.Step("Login",
			chromedp.WaitVisible(`#passwordInput`, chromedp.ByID),
			chromedp.WaitEnabled(`#passwordInput`, chromedp.ByID),
			chromedp.SendKeys(`#passwordInput`, s.User().Password()+"\r", chromedp.ByID),
		)
		_ = engine.StepTimeout("Wait loading", 5*time.Second, chromedp.WaitVisible(loadedID, chromedp.ByID))
		option := <-engine.Either("Login Check",
			szenario.EitherOption{ID: loadedID, Action: chromedp.WaitVisible(loadedID, chromedp.ByID)},
			szenario.EitherOption{ID: errorClass, Action: chromedp.WaitVisible(errorClass, chromedp.ByID)},
		)
		if option == loadedID {
			s.User().LoginSuccessfull()
			return nil
		}
		pw := s.User().NextPassword()
		if pw == "" {
			break
		}
	}
	return errors.New("Login failed")
}

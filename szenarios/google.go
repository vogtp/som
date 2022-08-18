package szenarios

import (
	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// GoogleSzenario does load google
type GoogleSzenario struct {
	*szenario.Base
}

// Execute the szenario
func (GoogleSzenario) Execute(engine szenario.Engine) (err error) {
	engine.Step("Loading",
		chromedp.Navigate("https://google.ch/"),
		chromedp.WaitVisible(`#tophf`, chromedp.ByID),
	)
	engine.Step("Check",
		engine.Body(engine.Contains("Google Search"), engine.Contains("About"), engine.Bigger(1000)),
	)

	return nil
}

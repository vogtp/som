package szenarios

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// GoogleSzenario does load google
type GoogleSzenario struct {
	*szenario.Base
	Search      string
	MustContain []string
}

// Execute the szenario
func (gs GoogleSzenario) Execute(engine szenario.Engine) (err error) {
	engine.Step("Loading",
		chromedp.Navigate("https://google.ch/"),
		chromedp.WaitVisible(`#tophf`, chromedp.ByID),
	)
	if len(gs.Search) > 1 {
		if ok := engine.IsPresent("#W0wltc", chromedp.ByID); ok {
			engine.Step("Click Accept", chromedp.Click("#W0wltc", chromedp.ByID))
		}
		searchInput := `document.querySelector("body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input")`
		engine.Step("searching",
			chromedp.WaitReady(searchInput, chromedp.ByJSPath),
			chromedp.SendKeys(
				searchInput,
				fmt.Sprintf("%s\r", gs.Search),
				chromedp.ByJSPath, // copy JSPath from chrom developer tools
			),
			chromedp.WaitVisible(`#result-stats`, chromedp.ByID),
		)
	}
	checks := make([]szenario.CheckFunc, 1, len(gs.MustContain)+1)
	checks[0] = engine.Bigger(1000)
	for _, c := range gs.MustContain {
		checks = append(checks, engine.Contains(c))
	}
	engine.Step("Check",
		engine.Body(checks...),
	)

	return nil
}

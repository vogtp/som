package szenarios

import (
	"context"
	"errors"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// CrasherSzenario a szenarion that loads google and then crashes
type CrasherSzenario struct {
	*szenario.Base
}

// Execute the szenario
func (CrasherSzenario) Execute(engine szenario.Engine) (err error) {
	engine.Step("Loading",
		chromedp.Navigate("https://google.ch/"),
		chromedp.WaitVisible(`#tophf`, chromedp.ByID),
	)
	engine.Step("Crash", chromedp.ActionFunc(func(ctx context.Context) error {
		engine.AddErr(errors.New("Test error 1"))
		engine.AddErr(errors.New("Test error 2"))
		engine.AddErr(errors.New("Test error 3"))
		return fmt.Errorf("requested to fail")
	}))

	return nil
}

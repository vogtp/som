package szenarios

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"image/png"

	"github.com/chromedp/chromedp"
	"github.com/nfnt/resize"
	"github.com/vogtp/som/pkg/monitor/cdp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// CrasherSzenario a szenarion that loads google and then crashes
type CrasherSzenario struct {
	*szenario.Base
}

// Execute the szenario
func (CrasherSzenario) Execute(engine szenario.Engine) (err error) {
	defer func() {
		// if there was an error dump the screenshot and the html code to std out
		r := recover()
		if err != nil || r != nil {
			if cdpE, ok := engine.(*cdp.Engine); ok {
				// print screenshot
				if sc := cdpE.ScreenShot(); len(sc) > 0 {
					printImage(sc, 500)
				}
				// print html
				if src := cdpE.GetHTML(); len(src) > 0 {
					fmt.Printf("<pre><code>%s</code></pre>\n", html.EscapeString(src))
				}
			}
			if r != nil {
				panic(r)
			}
		}
	}()
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

func printImage(img []byte, width int) {
	image, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return
	}

	newImage := resize.Resize(uint(width), 0, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	err = enc.Encode(buf, newImage)
	if err != nil {
		return
	}

	imgb64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	fmt.Printf("<img src='data:image/png;base64, %s' />\n", imgb64)
}

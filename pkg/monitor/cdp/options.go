package cdp

import (
	"time"
)

// Repeat sets the repeat intervall in munutes (default 0 - no looping)
func Repeat(r time.Duration) Option {
	return func(e *Engine) {
		if r > 0 {
			e.repeat = r
		}
	}
}

// Timeout sets the job timeout in seconds (default 60s)
func Timeout(t time.Duration) Option {
	return func(e *Engine) {
		if t > 0 {
			e.SetTimeout(t)
		}
	}
}

// NoClose prevents the bowser form closing (used for debugging)
func NoClose() Option {
	return func(e *Engine) {
		e.noClose = true
	}
}

// Show controlls if the browser is started headless or not
func Show() Option {
	return func(e *Engine) {
		e.show = true
	}
}

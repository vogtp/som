package cdp

import (
	"time"
)

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

// Headless controls if the browser is started headless or not
func Headless(b bool) Option {
	return func(e *Engine) {
		e.headless = b
	}
}

// StepBreakPoint sets a channel used a break point after each step
func StepBreakPoint(c chan any) Option {
	return func(e *Engine) {
		c <- "Start"
		e.stepBreakPoint = c
	}
}

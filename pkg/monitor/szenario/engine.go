package szenario

import (
	"time"

	"github.com/chromedp/chromedp"
)

// JobFunc is the signature of a function that runs as job
type JobFunc func(Engine) error

// CheckFunc is used to check the plaintext body content
type CheckFunc func(*string) error

// Engine is executing a szenario
type Engine interface {
	// StepTimeout executes a Step with an timeout
	StepTimeout(name string, timeout time.Duration, actions ...chromedp.Action) error
	// Step executes the actions given and records how long it takes
	Step(name string, actions ...chromedp.Action)
	// IsPresent checks if something is present
	IsPresent(sel interface{}, opts ...chromedp.QueryOption) bool
	// Body is used to check the content of the page
	Body(checks ...CheckFunc) chromedp.Action
	// WaitForEver blocks until the timeout is reached
	WaitForEver()
	// Headless indicates if the browser is headless (i.e. does not show on screen)
	Headless() bool
	// Dump prints the body and its size to log
	Dump() CheckFunc
	// Contains looks for a string in the body
	Contains(s string) CheckFunc
	// NotContains looks for a string in the body and errs if found
	NotContains(s string) CheckFunc
	// Bigger checks if the size of the body (in bytes) in bigger than i
	Bigger(i int) CheckFunc
	// Strings gets the body as plaintext
	Strings(html *string) CheckFunc
	// SetStatus sets a status of the event
	SetStatus(key, val string)
	// AddErr adds a error to the event
	AddErr(err error)
}

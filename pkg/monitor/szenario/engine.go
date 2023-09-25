package szenario

import (
	"time"

	"log/slog"

	"github.com/chromedp/chromedp"
)

// JobFunc is the signature of a function that runs as job
type JobFunc func(Engine) error

// CheckFunc is used to check the plaintext body content
type CheckFunc func(*string) error

// EitherOption options for either
type EitherOption struct {
	ID     any
	Action chromedp.Action
}

// Engine is executing a szenario
type Engine interface {
	// StepTimeout executes a Step with an timeout
	StepTimeout(name string, timeout time.Duration, actions ...chromedp.Action) error
	// Step executes the actions given and records how long it takes
	Step(name string, actions ...chromedp.Action)

	// SetInputField sets a HTML input field and validates that it has been set
	SetInputField(stepName string, sel interface{}, value string, opts ...func(*chromedp.Selector)) error
	// IsPresent checks if something is present
	IsPresent(sel interface{}, opts ...chromedp.QueryOption) bool
	// Either wait for a list of options and sends the name of the first met option to the channel
	Either(name string, option ...EitherOption) <-chan any
	// Body is used to check the content of the page
	Body(checks ...CheckFunc) chromedp.Action
	// IsHeadless indicates if the browser is headless (i.e. does not show on screen)
	IsHeadless() bool
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
	// Log returns the logger
	Log() *slog.Logger

	// WaitForEver blocks until the timeout is reached
	WaitForEver()
	// BreakWaitForUserInput waits until any key is clicked on the cmdlint
	BreakWaitForUserInput()
}

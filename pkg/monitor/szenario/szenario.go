package szenario

import "time"

// Szenario is the definition of a monitoring szenario
type Szenario interface {
	Name() string               // Name returns the name
	User() User                 // User returns the user the szenario runs with
	SetUser(u User)             // SetUser set the user the szenario runs with
	Execute(Engine) error       // Execute the szenario
	RepeatDelay() time.Duration // RepeatDelay between executions
	Timeout() time.Duration     // Timeout for an execution
	MaxLoginTry() int           // MaxLoginTry returns how many times a login with a new password should be attemped
}

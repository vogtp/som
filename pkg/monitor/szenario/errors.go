package szenario

import (
	"fmt"
	"time"
)

// TimeoutError indicated a timeout
type TimeoutError struct {
	Timeout time.Duration
}

// Error interface
func (te TimeoutError) Error() string {
	return fmt.Sprintf("%v timeout reached", te.Timeout)
}

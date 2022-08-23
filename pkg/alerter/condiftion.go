package alerter

import (
	"github.com/vogtp/som/pkg/core/msg"
)

// Conditon is a check if a alerting rule is triggered
type Conditon interface {
	Kind() string
	Check(*msg.AlertMsg) bool
}

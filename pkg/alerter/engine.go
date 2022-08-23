package alerter

import (
	"errors"

	"github.com/vogtp/som/pkg/core/msg"
)

// Engine sends out alerts
type Engine interface {
	Kind() string
	Send(*msg.AlertMsg, *Destination) error
}

var engines map[string]Engine = make(map[string]Engine)

// AddEngine add an engine (mail, chat etc) to alerting
func AddEngine(e Engine) error {
	if e == nil {
		return errors.New("engine must not be nil")
	}
	if len(e.Kind()) < 1 {
		return errors.New("engine must have a kind")
	}
	engines[e.Kind()] = e
	return nil
}

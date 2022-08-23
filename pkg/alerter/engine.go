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

// AddEngine add an engine (mail, chat etc) to alerting
func (a *Alerter) AddEngine(e Engine) error {
	if e == nil {
		return errors.New("engine must not be nil")
	}
	if len(e.Kind()) < 1 {
		return errors.New("engine must have a kind")
	}
	a.engines[e.Kind()] = e
	return nil
}

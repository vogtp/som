package gravbus

import (
	"github.com/suborbital/grav/grav"
)

type Message struct {
	grav.Message
}

func (m Message)ReplyTo() string{
	return m.UUID()
}
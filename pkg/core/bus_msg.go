package core

type BusMsgFuc func(BusMessage) error

type BusMessage interface {
	Data() []byte
	Type() string
	ReplyTo() string
}

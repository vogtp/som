package bus

type Message struct {
	RoutingKey string

	Payload []byte
}

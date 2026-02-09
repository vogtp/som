package bus

type Message struct {
	RoutingKey string

	Body []byte
}

package cfg

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
)

const (
	// BusURL The URL the bus backend should connect to
	BusURL = "bus.url"

	// BusTimeout Timeout for connections to the bus
	BusTimeout = "bus.timeout"
)

func busFlags() {
	pflag.String(BusURL, nats.DefaultURL, "The URL the NATS bus backend should connect to")
	pflag.Duration(BusTimeout, 15*time.Second, "Timeout for connections to the bus")
}

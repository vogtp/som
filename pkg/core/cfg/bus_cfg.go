package cfg

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
)

const (
	// BusURL The URL the bus backend should connect to
	BusURL = "bus.url"
	// BusMonitoringPort The port of the bus monitoring backend
	BusMonitoringPort = "bus.monitoring_port"
	// BusToken is the auth token for the bus
	BusToken = "bus.token"
	// BusTimeout Timeout for connections to the bus
	BusTimeout = "bus.timeout"
)

func busFlags() {
	pflag.String(BusURL, nats.DefaultURL, "The URL the NATS bus backend should connect to")
	pflag.Int(BusMonitoringPort, 8222, "The port of the bus monitoring backend")
	pflag.String(BusToken, "", "the auth token for the bus")
	pflag.Duration(BusTimeout, 15*time.Second, "Timeout for connections to the bus")
}

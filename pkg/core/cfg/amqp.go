package cfg

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	AmqpHost     = "amqp.host"
	AmqpPort     = "amqp.port"
	AmqpUser     = "amqp.user"
	AmqpPasswort = "amqp.password"
)

var amqpDefaultPassword = "" // set this in a init func in a ignored file in order for the tests to work

func amqpFlags() {
	pflag.String(AmqpHost, "localhost", "Host running the AMQP server (e.g. rabbitmq)")
	pflag.Int(AmqpPort, 5672, "Port the AMQP server is running on")
	pflag.String(AmqpUser, "som", "AMQP user")
	pflag.String(AmqpPasswort, amqpDefaultPassword, "AMQP user password")
}

func AmqpURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		viper.GetString(AmqpUser),
		viper.GetString(AmqpPasswort),
		viper.GetString(AmqpHost),
		viper.GetInt(AmqpPort),
	)
}

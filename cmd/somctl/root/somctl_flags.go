package root

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

const (
	// LogRawBus flag name to log the raw bus
	LogRawBus = "bus.log.raw"
	// StandAlone flag to start a local stater in oder to work standalone
	StandAlone = "standalone"
)

func init() {
	viper.SetDefault(cfg.LogLevel, defaultLogLevel.String())
	viper.SetDefault(cfg.CoreStartdelay, time.Millisecond)
}

func isCmdlineFlag(n string) bool {
	return strings.Contains(fmt.Sprintf("%v", os.Args), n)
}

func processFlags() {
	pflag.Bool(LogRawBus, false, "Log bus messages")
	pflag.Bool(StandAlone, false, "Run in standalone mode, i.e. start a stater in the background.")
	cfg.Parse()
	// only set loglevel from cmd line
	if !isCmdlineFlag(cfg.LogLevel) {
		viper.Set(cfg.LogLevel, defaultLogLevel)
	}
	if !isCmdlineFlag(cfg.CheckRepeat) {
		viper.Set(cfg.CheckRepeat, 0)
	}
}

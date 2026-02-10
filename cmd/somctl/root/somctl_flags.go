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
	// StandAlone flag to start a local stater in oder to work standalone
	StandAlone = "standalone"
)

func init() {
	viper.SetDefault(cfg.LogLevel, defaultLogLevel.String())
	viper.SetDefault(cfg.CoreStartdelay, time.Millisecond)
	rootCtl.PersistentFlags().Bool(StandAlone, false, "Run in standalone mode, i.e. start a stater in the background.")
	rootCtl.PersistentFlags().String(cfg.LogLevel, "warn", "Set the loglevel: error warn info debug trace off")
	rootCtl.PersistentFlags().Int(cfg.WebPort, 0, "Port the webserver runs on")
	rootCtl.PersistentFlags().Bool(cfg.LogSource, true, "Log the source line")
	rootCtl.PersistentFlags().Bool(cfg.LogJSON, false, "Log in json")
	rootCtl.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
}

func isCmdlineFlag(n string) bool {
	return strings.Contains(fmt.Sprintf("%v", os.Args), n)
}

func processFlags() {
	//	 cfg.Parse()
	// only set loglevel from cmd line
	if !isCmdlineFlag(cfg.LogLevel) {
		viper.Set(cfg.LogLevel, defaultLogLevel)
	}
	if !isCmdlineFlag(cfg.CheckRepeat) {
		viper.Set(cfg.CheckRepeat, 0)
	}
}

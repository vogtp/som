package check

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

var (
	defaultLogLevel = hclog.Error
)

func isCmdlineFlag(n string) bool {
	return strings.Contains(fmt.Sprintf("%v", os.Args), n)
}

func processFlags() {
	pflag.String(cfg.CheckUser, "", "User name of the user to run the check with")
	cfg.Parse()
	// only set loglevel from cmd line
	if !isCmdlineFlag(cfg.LogLevel) {
		viper.Set(cfg.LogLevel, defaultLogLevel)
	}

}

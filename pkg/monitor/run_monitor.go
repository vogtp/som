package monitor

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/monitor/cdp"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// Run the monitor
func Run(name string, coreOpts ...core.Option) (func(), error) {
	// if pflag.Lookup(cfg.CheckUser) == nil {
	// 	pflag.String(cfg.CheckUser, "", "User name of the user to run the check with")
	// }
	// pflag.Bool(cfg.BrowserShow, false, "Show the browser window")
	// pflag.Bool(cfg.BrowserNoClose, false, "Do not close the browser window in the end. Implies show, timeout 10m  and no repeat")
	// pflag.Duration(cfg.CheckTimeout, 60*time.Second, "Check timeout")
	// pflag.Duration(cfg.CheckRepeat, 0, "Check intervall (e.g. 5m)")
	// pflag.Duration(cfg.CheckStepDelay, 0, "Delay between steps (e.g. 100ms)")
	// pflag.String(cfg.CheckRegion, "default", "The region the check runs in")
	cfg.Parse()
	username := viper.GetString(cfg.CheckUser)
	if len(username) < 1 {
		return func() {}, fmt.Errorf("No user given. Use --%s or set it in the config", cfg.CheckUser)
	}
	c, close, err := core.New(fmt.Sprintf("%s.%s", name, username), coreOpts...)
	if err != nil {
		return nil, err
	}
	if c.SzenaioConfig() == szenario.NoConfig || c.SzenaioConfig().SzenarioCount() < 1 {
		panic("Monitor needs szenarios, no szenario config given")
	}
	go loop(c, username)

	return close, nil
}

func loop(c *core.Core, username string) {
	slog := c.Log().With(log.User, username)
	err := fmt.Errorf("Start")
	for err != nil {
		err = run(c, username)
		if err != nil {
			wait := 30 * time.Second
			slog.Error("Szenario run error", log.Error, err, "next_run_in", wait)
			time.Sleep(wait)
		}
	}
}

func run(c *core.Core, username string) error {
	cdp, cancel := cdp.New()
	defer cancel()
	err := cdp.RunUser(username)
	if err != nil {
		panic(err)
	}
	return nil
}

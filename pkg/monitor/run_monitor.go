package monitor

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/monitor/cdp"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

func init() {
	pflag.String(cfg.CheckUser, "", "User name of the user to run the check with")
}

// Run the monitor
func Run(name string, coreOpts ...core.Option) (func(), error) {
	cfg.Parse()
	username := viper.GetString(cfg.CheckUser)
	if len(username) < 1 && (hcl.IsGoRun() || hcl.IsGoTest()) {
		username = viper.GetString(cfg.CheckUser)
	}
	if len(username) < 1 {
		return func() {}, fmt.Errorf("No user given. Use --%s or set it in the config", cfg.CheckUser)
	}
	c, close := core.New(fmt.Sprintf("%s.user.%s", name, username), coreOpts...)
	if c.SzenaioConfig() == szenario.NoConfig || c.SzenaioConfig().SzenarioCount() < 1 {
		panic("Monitor needs szenarios, no szenario config given")
	}
	go loop(c, username)

	return close, nil
}

func loop(c *core.Core, username string) {
	hcl := c.HCL()
	err := fmt.Errorf("Start")
	for err != nil {
		err = run(c, username)
		if err != nil {
			hcl.Errorf("Szenario run: %v", err)
			wait := 30 * time.Second
			hcl.Errorf("Waiting %v", wait)
			<-time.After(wait)
		}
	}
}

func run(c *core.Core, username string) error {
	user, err := user.Store.Get(username)
	if err != nil {
		return fmt.Errorf("no such user: %v", username)
	}
	cdp, cancel := cdp.New()
	defer cancel()
	szs, err := c.SzenaioConfig().ByUser(user)
	if err != nil {
		panic(err)
	}
	cdp.Execute(szs...)
	return nil
}

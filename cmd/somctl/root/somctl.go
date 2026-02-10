package root

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/som/cmd/somctl/dbctl"
	"github.com/vogtp/som/cmd/somctl/incidentctl"
	"github.com/vogtp/som/cmd/somctl/szenarioctl"
	"github.com/vogtp/som/cmd/somctl/userctl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/bus"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater"
)

var (
	defaultLogLevel = hclog.Error
)

// AddCommand adds a *cobra.Command to somctl
func AddCommand(c *cobra.Command) {
	rootCtl.AddCommand(c)
}

// Command adds the root command
func Command(ctx context.Context, szCfg *szenario.Config) {
	processFlags()

	startCore(ctx, szCfg)

	rootCtl.AddCommand(userctl.Command())
	rootCtl.AddCommand(szenarioctl.Command())
	rootCtl.AddCommand(incidentctl.Command())
	rootCtl.AddCommand(dbctl.Command())

	if err := rootCtl.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
	}
}

func startCore(ctx context.Context, szCfg *szenario.Config) {
	if !viper.GetBool(StandAlone) {
		// normal mode: just start a core to connect to the mesh
		var err error
		c, coreClose, err = core.New("somctl", core.Szenario(szCfg))
		if err != nil {
			panic(err)
		}
	}
	//standalong mode: start a stater
	var err error
	coreClose, err = stater.Run(ctx, "somctl", core.Szenario(szCfg))
	if err != nil {
		fmt.Printf("Cannot start core: %v", err)
		os.Exit(-1)
	}
	c = core.Get()
}

var (
	c         *core.Core
	coreClose func()

	rootCtl = &cobra.Command{
		Use:   "somctl",
		Short: "Commandline interface to SOM",
		Long:  `Commandline interface to the Service Oriented Monitoring`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !cmd.IsAvailableCommand() {
				return
			}
			if viper.GetBool(LogRawBus) {
				c.Log().Info("Logging raw bus")
				c.Bus().Receive("*", func(subject string, msg *bus.Message) {
					fmt.Fprintf(cmd.OutOrStdout(), "Raw Bus: %s\n", string(msg.Body))
				})
			}
			time.Sleep(300 * time.Millisecond)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if !cmd.IsAvailableCommand() {
				return
			}
			core.Get().BusFIXME().WaitMsgProcessed()
			coreClose()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

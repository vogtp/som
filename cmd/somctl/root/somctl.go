package root

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/suborbital/grav/grav"
	"github.com/vogtp/som/cmd/somctl/dbctl"
	"github.com/vogtp/som/cmd/somctl/incidentctl"
	"github.com/vogtp/som/cmd/somctl/szenarioctl"
	"github.com/vogtp/som/cmd/somctl/userctl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

var (
	defaultLogLevel = hclog.Error
)

// AddCommand adds a *cobra.Command to somctl
func AddCommand(c *cobra.Command) {
	rootCtl.AddCommand(c)
}

// Command adds the root command
func Command(szCfg *szenario.Config) error {
	processFlags()

	c, coreClose = core.New("somctl", core.Szenario(szCfg))

	rootCtl.AddCommand(userctl.Command())
	rootCtl.AddCommand(szenarioctl.Command())
	rootCtl.AddCommand(incidentctl.Command())
	rootCtl.AddCommand(dbctl.Command())

	if err := rootCtl.Execute(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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
				c.Bus().Connect().On(func(m grav.Message) error {
					fmt.Fprintf(cmd.OutOrStdout(), "Raw Bus: %s\n", string(m.Data()))
					return nil
				})
			}
			<-time.After(300 * time.Millisecond)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if !cmd.IsAvailableCommand() {
				return
			}
			core.Get().Bus().WaitMsgProcessed()
			coreClose()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

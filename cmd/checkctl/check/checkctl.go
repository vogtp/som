package check

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
	"github.com/vogtp/som/cmd/checkctl/debugctl"
	"github.com/vogtp/som/cmd/checkctl/ldapctl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater"
)

// Command adds the root command
func Command(ctx context.Context, szCfg *szenario.Config) {
	processFlags()

	//startCore(szCfg)

	checkCtl.AddCommand(ldapctl.Command())
	checkCtl.AddCommand(debugctl.Command())
	director.GenerateDirectorConfigPFlag(checkCtl.PersistentFlags())
	checkCtl.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})

	if err := checkCtl.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
	}
}

var (
	c         *core.Core
	coreClose func()

	checkCtl = &cobra.Command{
		Use:   "checkctl",
		Short: "Run a check with a monitoring plugin interface",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if !cmd.IsAvailableCommand() {
				return
			}
			//FIXME core.Get().Bus().WaitMsgProcessed()
			// coreClose()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if director.ShouldGenerate() {
				d := director.Generator{
					NamePrefix:     "293",
					Description:    "CheckCtl: run Icinga2 commands",
					DescriptionURL: "https://github.com/vogtp/som/",
					CobraCmd:       cmd,
					Output:         os.Stdout,
					Criticality:    icinga.Criticality7x24,
				}
				if err := d.Generate(); err != nil {
					return err
				}
				os.Exit(0)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

// AddCommand adds a *cobra.Command to somctl
func AddCommand(c *cobra.Command) {
	checkCtl.AddCommand(c)
}

func startCore(szCfg *szenario.Config) {
	//standalong mode: start a stater
	var err error
	coreClose, err = stater.Run("checkctl", core.Szenario(szCfg))
	if err != nil {
		fmt.Printf("Cannot start core: %v", err)
		os.Exit(-1)
	}
	c = core.Get()
}

package szenarioctl

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

func init() {
	pflag.String(cfg.CheckUser, "", "User name of the user to run the check with")
}

var szenarioCtl = &cobra.Command{
	Use:     "szenario",
	Aliases: []string{"sz"},
	Short:   "Manage SOM szenarios",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Command initialises the szenario commands
func Command() *cobra.Command {
	szCfg := core.Get().SzenaioConfig()
	if szCfg != nil && szCfg != szenario.NoConfig {
		szenarioRun.ValidArgs = szCfg.GetUserTypes()
		szs, _ := szCfg.ByUser(&user.User{UserType: szenario.UserTypeAll})
		for _, sz := range szs {
			szenarioRun.ValidArgs = append(szenarioRun.ValidArgs, strings.ToLower(sz.Name()))
		}
		szenarioCtl.AddCommand(szenarioRun)
	}
	szenarioCtl.AddCommand(szenarioLog)
	return szenarioCtl
}

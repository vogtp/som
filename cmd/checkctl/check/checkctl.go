package check

import (
	"fmt"

	"github.com/spf13/cobra"
	goicinga "github.com/vogtp/go-icinga"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/som"
	"github.com/vogtp/som/cmd/checkctl/debugctl"
	"github.com/vogtp/som/cmd/checkctl/ldapctl"
	"github.com/vogtp/som/cmd/checkctl/szenarioapi"
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// Command adds the root command
func Command(szCfg *szenario.Config) {
	goicinga.VersionMajor = som.VersionMajor
	goicinga.VersionMinor = som.VersionMinor
	goicinga.VersionPatch = som.VersionPatch
	processFlags()

	//startCore(szCfg)

	checkCtl.AddCommand(ldapctl.Command())
	checkCtl.AddCommand(debugctl.Command())
	checkCtl.AddCommand(szenarioapi.Command())

	if err := checkCtl.Execute(); err != nil {
		fmt.Println(err)
	}
}

var (
	checkCtl = &check.Command{
		Use:            "checkctl",
		Short:          "Run a check with a monitoring plugin interface",
		NamePrefix:     "293",
		DescriptionURL: "https://github.com/vogtp/som/",
		Criticality:    icinga.Criticality7x24,

		Command: &cobra.Command{
			PersistentPostRun: func(cmd *cobra.Command, args []string) {
				if !cmd.IsAvailableCommand() {
					return
				}
				//FIXME core.Get().Bus().WaitMsgProcessed()
				// coreClose()
			},
		},
	}
)

// AddCommand adds a *cobra.Command to somctl
func AddCommand(c *cobra.Command) {
	checkCtl.AddCommand(c)
}

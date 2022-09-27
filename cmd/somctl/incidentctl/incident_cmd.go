package incidentctl

import (
	"github.com/spf13/cobra"
)

var incidentCtl = &cobra.Command{
	Use:     "incident",
	Aliases: []string{"incidents", "inci", "inc"},
	Short:   "Manage SOM incidents",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Command initialises the szenario commands
func Command() *cobra.Command {

	incidentCtl.AddCommand(incidentReplay)
	return incidentCtl
}

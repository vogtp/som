package incidentctl

import (
	"github.com/spf13/cobra"
)

var incidentCtl = &cobra.Command{
	Use:     "incident",
	Aliases: []string{"incidents", "inci", "inc"},
	Short:   "Manage SOM incidents",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Command initialises the szenario commands
func Command() *cobra.Command {

	incidentCtl.AddCommand(incidentReplay)
	return incidentCtl
}

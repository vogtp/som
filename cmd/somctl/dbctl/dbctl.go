package dbctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
)

// Command adds all user commands
func Command() *cobra.Command {
	dbCtl.AddCommand(dbThinOutCtl)
	return dbCtl
}

var dbCtl = &cobra.Command{
	Use:   "db",
	Short: "Manage SOM Database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dbThinOutCtl = &cobra.Command{
	Use:     "thinout",
	Short:   "Thin out entries of incidents to make free up space",
	Long:    ``,
	Aliases: []string{"thin"},
	RunE: func(cmd *cobra.Command, args []string) error {
		access, err := db.New()
		if err != nil {
			return fmt.Errorf("cannot connect to DB: %w", err)
		}
		return access.IncidentSummary.ThinOut(cmd.Context())
	},
}

package debugctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Command adds all ldap commands
func Command() *cobra.Command {
	debugCtl.AddCommand(debugEnvCtl)
	return debugCtl
}

var debugCtl = &cobra.Command{
	Use:   "debug",
	Short: "debug stuff",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var debugEnvCtl = &cobra.Command{
	Use:   "env",
	Short: "Show Environment",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Environment:\n")
		for _, e := range os.Environ() {
			fmt.Printf("  %v\n", e)
		}
		fmt.Printf("Arguments:\n")
		for i, a := range args {
			fmt.Printf("  %v: %v\n", i, a)
		}
		return nil
	},
}

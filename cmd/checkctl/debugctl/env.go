package debugctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Command adds all ldap commands
func Command() *cobra.Command {
	debugEnvCtl.AddCommand(debugEnvCMDCtl)
	debugCtl.AddCommand(debugEnvCtl)
	return debugCtl
}

var debugCtl = &cobra.Command{
	Use:   "debug",
	Short: "debug stuff",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// for _, f := range cmd.Flags() {

		// }
		pflag.VisitAll(func(f *pflag.Flag) {
			fmt.Printf("%+v\n", f)
		})
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

var debugEnvCMDCtl = &cobra.Command{
	Use:   "cmd",
	Short: "Show Environment",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%+v\n", getCmds(cmd))
		return nil
	},
}

func getCmds(cmd *cobra.Command) []string {
	n := cmd.Name()
	p := cmd.Parent()
	if p == nil {
		return []string{n}
	}
	cmds := getCmds(cmd.Parent())
	return append(cmds, n)
}

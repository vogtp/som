package szenarioctl

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
)

var szenarioList = &cobra.Command{
	Use:     "list",
	Short:   "list al SOM szenario",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		szConfig := core.Get().SzenaioConfig()
		szenarios, _ := szConfig.ByUser(&user.User{UserType: szenario.UserTypeAll})
		fmt.Println("User Types:")
		for _, n := range szConfig.GetUserTypes() {
			fmt.Printf("  %s %v\n", n, szConfig.GetUserType(n).Szenarios)
		}
		fmt.Println("\nSzenarios:")
		for _, s := range szenarios {
			fmt.Printf("  %s\n", strings.ToLower(s.Name()))
		}

		return nil
	},
}

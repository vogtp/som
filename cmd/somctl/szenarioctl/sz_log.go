package szenarioctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/msg"
)

var szenarioLog = &cobra.Command{
	Use:   "log",
	Short: "log all szenario bus messages",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		core.Get().Bus().Szenario.Handle(func(e *msg.SzenarioEvtMsg) {
			err := ""
			if e.Err() != nil {
				err = e.Err().Error()
			}
			fmt.Printf("%-20s %-10s %4.1f s %v\n", e.Name, e.Username, e.Counters["step.total"], err)
		})

		<-make(chan bool)
	},
}

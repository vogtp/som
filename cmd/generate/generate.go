package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var genCtl = &cobra.Command{
	Use:   "somgen",
	Short: "Generate needed SOM files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func main() {
	ctx, close := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer close()

	genCtl.AddCommand(genKeyCtl)

	if err := genCtl.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
	}
}

package szenarioapi

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/som"
)

const (
	somURL     = "som.url"
	szNameFlag = "szenario.name"
)

var (
	defaultSomURL       = ""
	defaultSzenarioName = "$host.name$"
)

// Command adds all SOM API commands
func Command() *cobra.Command {
	flags := somApiCtl.PersistentFlags()
	flags.String(somURL, defaultSomURL, "URL of the SOM Overview")
	flags.String(szNameFlag, defaultSzenarioName, "Szenarion nam,e")
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	return somApiCtl
}

var somApiCtl = &cobra.Command{
	Use:     "szenario",
	Short:   "Read a szenario status from the API",
	Long:    ``,
	Aliases: []string{"sz"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		check.SetWarningThresholdDefault("5s")
		check.SetCriticalThresholdDefault("10s")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return somApiResult(cmd, args)
	},
}

func somApiResult(cmd *cobra.Command, args []string) error {
	result := check.NewResult(cmd.Name(), check.CounterFormater(timeFormater))
	defer result.PrintExit()
	err := querySomAPI(cmd.Context(), result)

	if err != nil {
		result.SetError(err)
	}
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Value.String() == f.DefValue {
			return
		}
		result.SetStatus(f.Name, f.Value)
	})
	result.SetStatus("Version", som.Version)
	return err
}

func timeFormater(name string, value check.Data) string {
	t, ok := value.Value.(time.Duration)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf("%.2fs", t.Seconds())
}

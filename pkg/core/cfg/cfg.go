package cfg

import (
	"strings"
	"time"

	"log/slog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SetConfigFileName sets the config file name
func SetConfigFileName(n string) {
	viper.Set(CfgFile, n)
}

// Parse parses the config
func Parse() {
	if pflag.Parsed() {
		slog.Debug("pflags already parsed")
		return
	}

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		slog.Error("cannot bin flags", "error", err)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName(viper.GetString(CfgFile))
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/som/")
	viper.AddConfigPath("/som/")
	viper.AddConfigPath("$HOME/.som")
	viper.AddConfigPath("$SOM_HOME/")
	viper.AddConfigPath("$SOM_ROOT/")

	viper.SetEnvPrefix("SOM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	processConfigFile()
	if viper.GetBool(BrowserNoClose) {
		viper.Set(BrowserShow, true)
		viper.Set(CheckTimeout, 10*time.Minute)
		viper.Set(CheckRepeat, 0)
	}

}

func processConfigFile() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Debug("config not found", "error", err)
		} else {
			slog.Warn("config error", "error", err)
		}
	}
	if viper.GetBool(CfgSave) {
		//autosave the current config
		go func() {
			time.Sleep(10 * time.Second)
			for {
				if !viper.GetBool(CfgSave) {
					slog.Debug("Requested not to save config")
					return
				}
				// should we write it regular and overwrite?
				slog.Info("Writing config")
				if err := viper.WriteConfigAs(viper.GetString(CfgFile)); err != nil {
					slog.Warn("Could not write config", "error", err)
				}
				time.Sleep(time.Hour)
			}
		}()
	}
}

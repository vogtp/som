package cfg

import (
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
)

var (
	configFIleName = "som"
)

// SetConfigFileName sets the config file name
func SetConfigFileName(n string) {
	configFIleName = n
}

// Parse parses the config
func Parse() {
	if pflag.Parsed() {
		hcl.Debug("pflags already parsed")
		return
	}

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigType("yaml")
	viper.SetConfigName(configFIleName)
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
	if len(configFIleName) < 1 {
		return
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			hcl.Debugf("config not found: %v", err)
		} else {
			hcl.Warnf("config error: %v", err)
		}
	}
	//autosave the current config
	go func() {
		<-time.After(10 * time.Second)
		for {
			if !viper.GetBool(CfgSave) {
				hcl.Debug("Requested not to save config")
				return
			}
			// should we write it regular and overwrite?
			hcl.Info("Writing config")
			if err := viper.WriteConfigAs(configFIleName); err != nil {
				hcl.Warnf("Could not write config: %v", err)
			}
			<-time.After(time.Hour)
		}
	}()
}

// HclOptions configures the HCL logger (using commandline flags)
func HclOptions() hcl.LoggerOpt {
	Parse()
	logLvl := viper.GetString(LogLevel)
	if logLvl != "" {
		lvl := hclog.LevelFromString(logLvl)
		if lvl == hclog.NoLevel {
			hcl.Errorf("Unrecoginsed loglevel: %s", logLvl)
		} else {
			return hcl.WithLevel(lvl)
		}
	}
	return func(l *hcl.Logger) {}
}

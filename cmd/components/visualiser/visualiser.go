package main

import (
	"github.com/spf13/pflag"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/visualiser"
	"github.com/vogtp/som/szenarios"
)

func main() {
	// szenarios.Load() has to be replace by ones own szenario config
	szCfg := szenarios.Load()
	pflag.String(cfg.DataDir, "data", "Folder to save output like screenshots in")
	close, err := visualiser.Run("som.visualiser", core.Szenario(szCfg))
	defer close()
	if err != nil {
		panic(err)
	}
	// wait for ever
	<-make(chan any)
}

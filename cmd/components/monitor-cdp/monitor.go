package main

import (
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/monitor"
	"github.com/vogtp/som/szenarios"
)

func main() {
	// szenarios.Load() has to be replace by ones own szenario config
	szCfg := szenarios.Load()
	close, err := monitor.Run("som.monitor-cdp", core.Szenario(szCfg))
	defer close()
	if err != nil {
		panic(err)
	}
	// wait for ever
	<-make(chan any)
}

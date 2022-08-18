package main

import (
	"github.com/vogtp/som/pkg/alerter"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/monitor"
	"github.com/vogtp/som/pkg/stater"
	"github.com/vogtp/som/pkg/visualiser"
	"github.com/vogtp/som/szenarios"
)

func main() {
	// szenarios.Load() has to be replace by ones own szenario config
	szCfg := szenarios.Load()
	opts := []core.Option{
		core.Szenario(szCfg),
		core.WebPort(8083),
	}
	name := "som.allinone"
	close, err := stater.Run(name, opts...)
	defer close()
	if err != nil {
		panic(err)
	}
	if _, err := visualiser.Run(name, opts...); err != nil {
		panic(err)
	}
	if _, err := alerter.Run(name, opts...); err != nil {
		panic(err)
	}
	if _, err := monitor.Run(name, opts...); err != nil {
		panic(err)
	}
	// wait for ever
	<-make(chan any)
}

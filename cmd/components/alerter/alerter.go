package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vogtp/som/pkg/alerter"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/szenarios"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	// szenarios.Load() has to be replace by ones own szenario config
	szCfg := szenarios.Load()
	close, err := alerter.Run("som.alerter", core.Szenario(szCfg))
	defer close()
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
}

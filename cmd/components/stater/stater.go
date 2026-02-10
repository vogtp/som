package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/stater"
	"github.com/vogtp/som/szenarios"
)

func main() {
	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer close()
	// szenarios.Load() has to be replace by ones own szenario config
	szCfg := szenarios.Load()
	close, err := stater.Run(ctx, "som.stater", core.Szenario(szCfg))
	defer close()
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
}

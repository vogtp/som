package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vogtp/som/cmd/somctl/root"
	"github.com/vogtp/som/szenarios"
)

func main() {
	ctx, stop:=signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	// szenarios.Load() has to be replace by ones own szenario config
	root.Command(ctx, szenarios.Load())
}


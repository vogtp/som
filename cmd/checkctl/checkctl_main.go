package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vogtp/som/cmd/checkctl/check"
	"github.com/vogtp/som/szenarios"
)

func main() {
	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer close()
	check.Command(ctx, szenarios.Load())
}

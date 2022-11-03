package main

import (
	"github.com/vogtp/som/cmd/somctl/root"
	"github.com/vogtp/som/szenarios"
)

func main() {
	// szenarios.Load() has to be replace by ones own szenario config
	root.Command(szenarios.Load())
}

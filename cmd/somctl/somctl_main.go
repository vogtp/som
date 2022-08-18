package main

import (
	"github.com/vogtp/som/cmd/somctl/root"
	"github.com/vogtp/som/szenarios"
)

func main() {
	// szenarios.Load() has to be replace by ones own szenario config
	szCfg := szenarios.Load()
	root.Command(szCfg)
}

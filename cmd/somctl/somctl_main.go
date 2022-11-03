package main

import (
	"fmt"

	"github.com/vogtp/som/cmd/somctl/root"
	"github.com/vogtp/som/szenarios"
)

func main() {
	// szenarios.Load() has to be replace by ones own szenario config
	if err := root.Command(szenarios.Load()); err != nil {
		fmt.Printf("Root command failed: %v", err)
	}
}

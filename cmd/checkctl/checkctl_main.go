package main

import (
	"github.com/vogtp/som/cmd/checkctl/check"
	"github.com/vogtp/som/szenarios"
)

func main() {
	check.Command(szenarios.Load())
}

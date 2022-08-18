package main

import (
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	t.Skip("Only for debugging")
	cmdLineArgs := "--check.user=vogtpa szenario run all "

	os.Args = append(os.Args, strings.Split(cmdLineArgs, " ")...)
	main()
}

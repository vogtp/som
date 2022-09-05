package som

import (
	_ "embed" // embed needs it
)

var (
	//go:embed README.md
	README []byte
)

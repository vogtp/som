package som

import (
	_ "embed" // embed needs it
)

var (
	//go:embed README.md
	// README to display on the visualiser
	README []byte
)

package env

import (
	"strings"

	"github.com/vogtp/som"
)

func GitDev() bool {
	return strings.EqualFold(som.Branch, "dev")
}

func GitQM() bool {
	return strings.EqualFold(som.Branch, "qm")
}

func GitProd() bool {
	return strings.EqualFold(som.Branch, "main") || strings.EqualFold(som.Branch, "master")
}

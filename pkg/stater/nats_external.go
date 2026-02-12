//go:build !embedd_nats

package stater

import (
	"github.com/vogtp/som/pkg/core"
)

func startEmbeddedNats(_ *core.Core) error {
	return nil
}

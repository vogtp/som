//go:build !embedd_nats

package stater

import (
	"github.com/vogtp/som/pkg/core"
)

func startEmbeddedNats(core *core.Core) error {
	return nil
}

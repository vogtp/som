package ent

import "github.com/vogtp/som/pkg/core/status"

// Level returns the level
func (a Alert) Level() status.Level {
	return status.Level(a.IntLevel)
}

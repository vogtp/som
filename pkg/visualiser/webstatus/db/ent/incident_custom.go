package ent

import "github.com/vogtp/som/pkg/core/status"

// Level returns the level
func (i Incident) Level() status.Level {
	return status.Level(i.IntLevel)
}

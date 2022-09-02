package status

import (
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

// Cleanup removes groups with Level below OK (i.e. unknown)
// and no events since status.cleanup config
func Cleanup(s Status) {
	cleanupDuration := viper.GetDuration(cfg.StatusCleanup)
	grp, ok := s.(*statusGroup)
	if !ok {
		return
	}
	children := make([]Grouper, 0, len(grp.children))
	for _, c := range grp.Szenarios() {
		if c.Level() < OK && time.Since(c.LastUpdate()) > cleanupDuration {
			continue
		}
		clenupSz(c, cleanupDuration)
		children = append(children, c)
	}
	grp.children = children
}

func clenupSz(grpr SzenarioGroup, cleanupDuration time.Duration) {
	grp, ok := grpr.(*szGroup)
	if !ok {
		return
	}
	children := make([]Grouper, 0, len(grp.children))
	for _, c := range grp.Regions() {
		if c.Level() < OK && time.Since(c.LastUpdate()) > cleanupDuration {
			continue
		}
		clenupRg(c, cleanupDuration)
		children = append(children, c)
	}
	grp.children = children
}

func clenupRg(grpr RegionGroup, cleanupDuration time.Duration) {
	grp, ok := grpr.(*regGroup)
	if !ok {
		return
	}
	children := make([]Grouper, 0, len(grp.children))
	for _, c := range grp.Users() {
		if c.Level() < OK && time.Since(c.LastUpdate()) > cleanupDuration {
			continue
		}
		children = append(children, c)
	}
	grp.children = children
}

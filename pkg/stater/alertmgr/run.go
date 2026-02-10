package alertmgr

import "log/slog"

// Run starts a alert manager
func Run(slog *slog.Logger) error {
	New(slog)
	return nil
}

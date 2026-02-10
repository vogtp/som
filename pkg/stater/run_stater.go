package stater

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/stater/alertmgr"
	"github.com/vogtp/som/pkg/stater/user"
)

// Run the stater
func Run(ctx context.Context, name string, coreOpts ...core.Option) (func(), error) {
	viper.Set(cfg.CoreStartdelay, 100*time.Millisecond)
	cfg.Parse()
	c, close, err := core.New(name, coreOpts...)
	if err != nil {
		return nil, err
	}
	user.IntialiseStore(ctx)

	if err := alertmgr.Run(); err != nil {
		c.Log().Warn("alertmgr refused to run", log.Error, err)
		return close, err
	}

	return close, nil
}

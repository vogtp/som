package stater

import (
	"context"
	"fmt"
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
	if err := startEmbeddedNats(c); err != nil {
		c.Log().Error("Cannot find a NATS bus nor start a embedded one", log.Error, err)
		close()
		return nil, fmt.Errorf("start embedded core: %w", err)
	}
	user.IntialiseStore(ctx)

	if err := alertmgr.Run(c.Log()); err != nil {
		c.Log().Warn("alertmgr refused to run", log.Error, err)
		return close, err
	}

	return close, nil
}

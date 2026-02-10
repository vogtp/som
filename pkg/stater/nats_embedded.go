package stater

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
)

func isNatsRunning(core *core.Core) (ret bool) {
	defer func() {
		if recover() != nil {
			ret = false
		}
	}()
	if err := core.Bus().EnsureConnected(); err == nil {
		return true
	}
	return false
}

func startEmbeddedNats(core *core.Core) error {
	sl := core.Log().With(log.Component, "embedded.nats.start")
	if isNatsRunning(core) {
		return nil
	}

	ns, err := server.NewServer(natsOptions(sl))

	if err != nil {
		return fmt.Errorf("starting embedded server: %w", err)
	}
	go ns.Start()

	if !ns.ReadyForConnections(viper.GetDuration(cfg.BusTimeout)) {
		return fmt.Errorf("connection timedout: %w", err)
	}
	sl.Info("Started embedded core", "bus.url", ns.ClientURL())
	return nil
}

func natsOptions(slog *slog.Logger) *server.Options {
	opts := &server.Options{
	}
	s := viper.GetString(cfg.BusURL)
	url, err := url.Parse(s)
	if err == nil {
		hs := strings.Split(url.Host, ":")
		if len(hs) > 1 {
			i, err := strconv.Atoi(hs[1])
			if err == nil {
				opts.Host = hs[0]
				opts.Port = i
			}
		}
		slog.Debug("Configuring NATS server with URL", "url_raw", s, "url", url, "nats_options", opts)
	}
	return opts
}

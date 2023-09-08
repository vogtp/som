package log

import (
	"io"
	"os"
	"strings"

	"log/slog"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core/cfg"
)

func New(name string) *slog.Logger {
	lvl := LevelFromString(viper.GetString(cfg.LogLevel))
	return Create(name, lvl)
}

func Create(name string, lvl slog.Level) *slog.Logger {
	logOpts := slog.HandlerOptions{
		Level: lvl,
	}
	logOpts.AddSource = viper.GetBool(cfg.LogSource)
	logJson := viper.GetBool(cfg.LogJson)
	if logOpts.AddSource {
		logOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey && len(groups) == 0 {
				return ProcessSourceField(a, logJson)
			}
			return a
		}
	}
	slog.Default()
	var logWriter io.Writer
	logWriter = os.Stdout
	var handler slog.Handler
	handler = slog.NewTextHandler(logWriter, &logOpts)
	if logJson {
		handler = slog.NewJSONHandler(logWriter, &logOpts)
	}
	log := slog.New(handler)
	log = log.With(slog.String("app", name))

	return log
}

func LevelFromString(levelStr string) slog.Level {
	// We don't care about case. Accept both "INFO" and "info".
	levelStr = strings.ToLower(strings.TrimSpace(levelStr))
	switch levelStr {
	case "trace":
		return slog.LevelDebug
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "off":
		return slog.Level(88)
	default:
		return slog.LevelInfo
	}
}

package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
)

func Setup(cfg config.LoggingConfig) {
	lv := parseLevel(cfg.Level)
	format := cfg.Format
	if format == "" {
		format = "console"
	}

	switch format {
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lv,
		})))
	default:
		slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      lv,
			TimeFormat: time.RFC3339,
			AddSource:  lv == slog.LevelDebug,
		})))
	}
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

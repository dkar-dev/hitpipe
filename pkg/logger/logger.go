package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

const (
	envLocal   = "local"
	envDev     = "dev"
	envStaging = "staging"
	envProd    = "prod"
)

func NewLogger(env string, level string) *slog.Logger {

	var log *slog.Logger

	var logLevel slog.Level

	switch level {

	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError

	}

	switch env {

	case envLocal:
		log = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel, TimeFormat: time.Kitchen}))

	case envDev:
		log = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel, TimeFormat: time.Kitchen}))

	case envStaging, envProd:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))
	}

	return log
}

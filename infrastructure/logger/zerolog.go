package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Config struct {
	Level       string
	ServiceName string
	Environment string
}

func init() {
	Configure(Config{
		Level:       os.Getenv("LOG_LEVEL"),
		ServiceName: os.Getenv("OTEL_SERVICE_NAME"),
		Environment: os.Getenv("OTEL_ENVIRONMENT"),
	})
}

// Configure sets the global structured logger.
func Configure(cfg Config) {
	level := zerolog.InfoLevel
	if cfg.Level != "" {
		parsedLevel, err := zerolog.ParseLevel(cfg.Level)
		if err == nil {
			level = parsedLevel
		}
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	context := zerolog.New(os.Stdout).With().Timestamp().Caller()
	if cfg.ServiceName != "" {
		context = context.Str("service", cfg.ServiceName)
	}
	if cfg.Environment != "" {
		context = context.Str("environment", cfg.Environment)
	}

	log.Logger = context.Logger()
}

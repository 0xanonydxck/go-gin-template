package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

type Mode string

func (m Mode) String() string {
	return string(m)
}

const (
	ProductionMode  = Mode("production")
	DevelopmentMode = Mode("development")
)

var (
	MODE       Mode
	PORT       int
	LIMIT_RATE string
	LOG_LEVEL  string

	OTEL_ENABLED                bool
	OTEL_SERVICE_NAME           string
	OTEL_ENVIRONMENT            string
	OTEL_EXPORTER_OTLP_ENDPOINT string
	OTEL_SAMPLE_RATIO           float64

	CORS_ALLOWED_ORIGINS string
	CORS_ALLOWED_METHODS string
	CORS_ALLOWED_HEADERS string
	CORS_EXPOSED_HEADERS string
	CORS_MAX_AGE         int

	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_PASSWORD string
	REDIS_DB       int

	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string

	ACCESS_SECRET  string
	REFRESH_SECRET string
)

func Init() {
	MODE = ModeEnv("MODE")
	LIMIT_RATE = StringEnv("LIMIT_RATE")
	PORT = IntEnv("PORT")
	LOG_LEVEL = StringEnvDefault("LOG_LEVEL", "info")

	OTEL_ENABLED = BoolEnvDefault("OTEL_ENABLED", false)
	OTEL_SERVICE_NAME = StringEnvDefault("OTEL_SERVICE_NAME", "simple-bookstore")
	OTEL_ENVIRONMENT = StringEnvDefault("OTEL_ENVIRONMENT", MODE.String())
	OTEL_EXPORTER_OTLP_ENDPOINT = StringEnvDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	OTEL_SAMPLE_RATIO = FloatEnvDefault("OTEL_SAMPLE_RATIO", 1)

	CORS_ALLOWED_ORIGINS = StringEnv("CORS_ALLOWED_ORIGINS")
	CORS_ALLOWED_METHODS = StringEnv("CORS_ALLOWED_METHODS")
	CORS_ALLOWED_HEADERS = StringEnv("CORS_ALLOWED_HEADERS")
	CORS_EXPOSED_HEADERS = StringEnv("CORS_EXPOSED_HEADERS")
	CORS_MAX_AGE = IntEnv("CORS_MAX_AGE")

	REDIS_HOST = StringEnv("REDIS_HOST")
	REDIS_PORT = StringEnv("REDIS_PORT")
	REDIS_PASSWORD = StringEnv("REDIS_PASSWORD")
	REDIS_DB = IntEnv("REDIS_DB")

	POSTGRES_HOST = StringEnv("POSTGRES_HOST")
	POSTGRES_PORT = StringEnv("POSTGRES_PORT")
	POSTGRES_USER = StringEnv("POSTGRES_USER")
	POSTGRES_PASSWORD = StringEnv("POSTGRES_PASSWORD")
	POSTGRES_DB = StringEnv("POSTGRES_DB")

	ACCESS_SECRET = StringEnv("ACCESS_SECRET")
	REFRESH_SECRET = StringEnv("REFRESH_SECRET")
}

func ModeEnv(key string) Mode {
	switch os.Getenv(key) {
	case ProductionMode.String():
		return ProductionMode
	case DevelopmentMode.String():
		return DevelopmentMode
	default:
		log.Fatal().Msg("🚨 mode must be prod or dev")
		return ""
	}
}

func IntEnv(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert " + key + " to int")
	}
	return value
}

func StringEnv(key string) string {
	return os.Getenv(key)
}

func StringEnvDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func BoolEnvDefault(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert " + key + " to bool")
	}
	return parsed
}

func FloatEnvDefault(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert " + key + " to float")
	}
	return parsed
}

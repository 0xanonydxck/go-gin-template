package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

var (
	PORT       int
	LIMIT_RATE string

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
	var err error

	LIMIT_RATE = os.Getenv("LIMIT_RATE")
	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert PORT to int")
	}

	CORS_ALLOWED_ORIGINS = os.Getenv("CORS_ALLOWED_ORIGINS")
	CORS_ALLOWED_METHODS = os.Getenv("CORS_ALLOWED_METHODS")
	CORS_ALLOWED_HEADERS = os.Getenv("CORS_ALLOWED_HEADERS")
	CORS_EXPOSED_HEADERS = os.Getenv("CORS_EXPOSED_HEADERS")
	CORS_MAX_AGE, err = strconv.Atoi(os.Getenv("CORS_MAX_AGE"))
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert CORS_MAX_AGE to int")
	}

	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	REDIS_DB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert REDIS_DB to int")
	}

	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB = os.Getenv("POSTGRES_DB")

	ACCESS_SECRET = os.Getenv("ACCESS_SECRET")
	REFRESH_SECRET = os.Getenv("REFRESH_SECRET")
}

package db

import (
	"context"
	"fmt"

	"github.com/chai-rs/simple-bookstore/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Redis returns the Redis client.
func Redis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT),
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})

	if config.OTEL_ENABLED {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			log.Warn().Err(err).Msg("failed to enable redis tracing")
		}
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("💣 failed to connect to redis")
	}

	return rdb
}

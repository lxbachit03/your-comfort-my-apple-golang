package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	goredis "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func NewRedisClient() *goredis.Client {

	redisAddress := fmt.Sprintf("%s:%d", config.AppConfig.Cache.Redis.Host, config.AppConfig.Cache.Redis.Port)

	cfg := RedisConfig{
		Addr:     redisAddress,
		Username: config.AppConfig.Cache.Redis.User,
		Password: config.AppConfig.Cache.Redis.Password,
		DB:       config.AppConfig.Cache.Redis.DB,
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	logger.Log.Info().Msg("🍺 Connected Redis")

	return client
}

package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	redisConnection *redis.Client
)

func GetConnection() (*redis.Client, error) {
	if redisConnection != nil {
		if _, err := redisConnection.Ping(context.Background()).Result(); err != nil {
			redisConnection = nil
			return nil, errors.New("cannot get connection")
		}
	}

	if redisConnection == nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.database"),
		})

		redisConnection = rdb
	}

	return redisConnection, nil
}

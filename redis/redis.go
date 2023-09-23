package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	redis *redis.Client
}

func NewRedis(addr string, password string, db int) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Redis{
		redis: rdb,
	}
}

func (r *Redis) PSub(ctx context.Context, channels string) <-chan *redis.Message {
	return r.redis.PSubscribe(ctx, channels).Channel()
}

package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	// channel
	REVOKE          = "revoke"
	MALICIOUS_ALTER = "maliciousAlter"
	ETH_BUY_SIGNAL  = "ethBuySignal"

	// accounts keys
	ACCOUNT_INFO  = "account"
	UUID          = "uuid"
	OTP_2_CHAT_ID = "otp2chatid"
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
	cmd := rdb.Ping(context.Background())
	if err := cmd.Err(); err != nil {
		logrus.WithError(err).Fatal("redis ping")
	}
	logrus.Info("connected redis")

	return &Redis{
		redis: rdb,
	}
}

func (r *Redis) PSub(ctx context.Context, channels string) <-chan *redis.Message {
	return r.redis.PSubscribe(ctx, channels).Channel()
}

func (r *Redis) Publish(ctx context.Context, channel string, message interface{}) error {
	cmd := r.redis.Publish(ctx, channel, message)
	if err := cmd.Err(); err != nil {
		logrus.WithError(err).Errorf("Publish: channel: %s, message: %s", channel, message)
		return err
	}
	return nil
}

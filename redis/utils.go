package redis

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (r *Redis) setAndCheck(ctx context.Context, key string, value interface{}, logMsg string) (err error) {
	statusCmd := r.redis.Set(ctx, key, value, 0)

	_, err = statusCmd.Result()
	if err != nil {
		logrus.WithError(err).Error(logMsg)
		return
	}
	return
}

package worker

import (
	"context"
	// "time"

	"github.com/go-redis/redis/v8"
)

func (rtp *RedisTaskProcessor) Get(ctx context.Context, key string) (string, error) {
	// Retrieve data from Redis using the provided key.
	val, err := rtp.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (rtp *RedisTaskProcessor) Set(key string, value string) error {
	// Set the value in Redis for the provided key.
	err := rtp.redis.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
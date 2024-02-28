package redis

import (
	"context"
	"time"
)

func SetKeyValue(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetKeyValue(ctx context.Context, key string) (string, error) {
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

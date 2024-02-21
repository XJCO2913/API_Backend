package zlog

import (
	"context"

	"api.backend.xjco2913/dao/redis"
)

const (
	LOG_KEY = "LOG"
)

type RedisWriter struct{}

// Write() implements the io.writer interface
func (r *RedisWriter) Write(logMsg []byte) (int, error) {
	ctx := context.Background()
	err := redis.RDB().LPush(ctx, LOG_KEY, logMsg).Err()
	if err != nil {
		return 0, err
	}

	return len(logMsg), nil
}
package main

import (
	"context"

	"api.backend.xjco2913/dao/redis"
	"api.backend.xjco2913/util/zlog"
)

const (
	LOG_KEY     = "LOG"
	ERR_LOG_KEY = "ERR_LOG"
)

type (
	RedisWriter    struct{}
	RedisErrWriter struct{}
)

// Write() implements the io.writer interface
func (r *RedisWriter) Write(logMsg []byte) (int, error) {
	ctx := context.Background()
	err := redis.RDB().RPush(ctx, LOG_KEY, logMsg).Err()
	if err != nil {
		return 0, err
	}

	return len(logMsg), nil
}

// Write() implements the io.writer interface
func (r *RedisErrWriter) Write(errMsg []byte) (int, error) {
	ctx := context.Background()
	err := redis.RDB().RPush(ctx, ERR_LOG_KEY, errMsg).Err()
	if err != nil {
		return 0, err
	}

	return len(errMsg), nil
}

var (
	localRedisLogWriter = &RedisWriter{}
	localRedisErrWriter = &RedisErrWriter{}
)

func init() {
	// logger dependency injection
	zlog.WithNewWriter(localRedisLogWriter, false)
	zlog.WithNewWriter(localRedisErrWriter, true)
}

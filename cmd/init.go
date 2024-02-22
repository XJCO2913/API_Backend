package main

import (
	"context"
	"log"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/redis"
	"api.backend.xjco2913/util/zlog"
	"gorm.io/gorm/logger"
)

const (
	LOG_KEY     = "LOG"
	ERR_LOG_KEY = "ERR_LOG"
	SQL_LOG_KEY = "SQL_LOG"
)

type (
	RedisWriter    struct{}
	RedisErrWriter struct{}
	RedisSQLWrirer struct{}
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

func (r *RedisSQLWrirer) Write(sqlMsg []byte) (int, error) {
	ctx := context.Background()
	err := redis.RDB().RPush(ctx, SQL_LOG_KEY, sqlMsg).Err()
	if err != nil {
		return 0, err
	}

	return len(sqlMsg), nil
}

var (
	localRedisLogWriter = &RedisWriter{}
	localRedisErrWriter = &RedisErrWriter{}
	localRedisSQLWriter = &RedisSQLWrirer{}
)

func init() {
	// logger dependency injection
	// server log
	zlog.WithNewWriter(localRedisLogWriter, false)
	zlog.WithNewWriter(localRedisErrWriter, true)
	// sql log
	dao.DB.Logger = logger.New(
		log.New(localRedisSQLWriter, "", log.LstdFlags),
		logger.Config{
			SlowThreshold: 800 * time.Millisecond,
			LogLevel: logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful: false,
		},
	)
}

package main

import (
	"context"
	"time"

	"api.backend.xjco2913/dao/redis"
)

func FlushLogs(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		redis.SyncLogs(ctx, LOG_KEY)
		redis.SyncLogs(ctx, ERR_LOG_KEY)
		redis.SyncLogs(ctx, SQL_LOG_KEY)
	}
}

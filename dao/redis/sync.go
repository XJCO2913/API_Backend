package redis

import (
	"context"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

// SyncLogs() will run async
func SyncLogs(ctx context.Context, logKey string) {
	// get the length of log list
	length, err := RDB().LLen(ctx, logKey).Result()
	if err != nil {
		return
	}

	if length <= 0 {
		// if no log, return
		return
	}

	logs, err := RDB().LRange(ctx, logKey, 0, length-1).Result()
	if err != nil {
		return
	}

	logDao := query.Use(dao.DB).Log
	for _, log := range logs {
		err := logDao.WithContext(ctx).Create(&model.Log{
			Msg:  log,
			Type: logKey,
		})
		if err != nil {
			return
		}
	}

	err = RDB().LTrim(ctx, logKey, length+1, -1).Err()
	if err != nil {
		return
	}
}
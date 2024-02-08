package zlog_test

import (
	"testing"

	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	zlog.Debug("debug msg", zap.String("type", "debug"))
	zlog.Warn("warn msg", zap.String("type", "warn"))
	zlog.Error("error msg", zap.Bool("error", true))
}
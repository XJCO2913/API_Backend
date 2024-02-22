package zlog

import (
	"io"
	"os"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	localLogger  *zap.Logger
	localCore    zapcore.Core    // default core
	localEncoder zapcore.Encoder // default encoder
)

func init() {
	// set time encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	localEncoder = zapcore.NewJSONEncoder(encoderConfig)

	// print into both stdout and log file
	localCore = zapcore.NewTee(
		zapcore.NewCore(localEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		// zapcore.NewCore(encoder, zapcore.AddSync(redisLogger), zapcore.DebugLevel),
	)

	localLogger = zap.New(localCore)
}

// Injection a new writer into zap logger
func WithNewWriter(w io.Writer, isError bool) {
	if isError {
		localCore = zapcore.NewTee(
			localCore,
			zapcore.NewCore(localEncoder, zapcore.AddSync(w), zapcore.ErrorLevel),
		)
	} else {
		localCore = zapcore.NewTee(
			localCore,
			zapcore.NewCore(localEncoder, zapcore.AddSync(w), zapcore.DebugLevel),
		)
	}

	localLogger = localLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return localCore
	}))
}

func Info(msg string, fields ...zap.Field) {
	callerFields := getCaller()
	fields = append(fields, callerFields...)
	localLogger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	callFields := getCaller()
	fields = append(fields, callFields...)
	localLogger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	callFields := getCaller()
	fields = append(fields, callFields...)
	localLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	callFields := getCaller()
	fields = append(fields, callFields...)
	localLogger.Error(msg, fields...)
}

func getCaller() []zap.Field {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return []zap.Field{}
	}

	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) // trim, only keep the func name

	callerFields := []zap.Field{}
	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return callerFields
}

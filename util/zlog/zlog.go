package zlog

import (
	"os"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	localLogger *zap.Logger
)

func init() {
	// set time encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	file, _ := os.OpenFile("./log/backend.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 644)
	fileWriterSyncer := zapcore.AddSync(file)

	// print into both stdout and log file
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriterSyncer, zapcore.DebugLevel),
	)

	localLogger = zap.New(core)
}

func Info(msg string, fields ...zap.Field) {
	callerFields := getCaller()
	fields = append(fields, callerFields...)
	localLogger.Info(msg, fields...)
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

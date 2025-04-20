package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", "./logs/server.log"}

	logger, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	Log = logger
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

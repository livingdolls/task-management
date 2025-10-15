package logger

import (
	"go.uber.org/zap"
)

var log, _ = zap.NewDevelopment() // atau zap.NewProduction() untuk JSON

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
	_ = log.Sync()
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
	_ = log.Sync()
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
	_ = log.Sync()
}

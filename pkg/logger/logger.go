package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(env string) {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func Fatal(msg string, fields ...zap.Field) { Log.Fatal(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Log.Error(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { Log.Info(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { Log.Debug(msg, fields...) }

package xlog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the global zap logger based on the APP_ENV environment variable.
// In "prod" mode it uses JSON encoding at InfoLevel; otherwise it uses console encoding with colors at DebugLevel.
func Init() error {
	env := os.Getenv("APP_ENV")

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		CallerKey:  "caller",
		MessageKey: "msg",

		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	var level zapcore.Level

	if env == "prod" {
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
		level = zapcore.InfoLevel
	} else {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return nil
}

// Sync flushes any buffered log entries.
func Sync() {
	_ = log.Sync()
}

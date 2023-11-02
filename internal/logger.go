package internal

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(level string) *zap.SugaredLogger {
	levels := map[string]zap.AtomicLevel{
		"debug": zap.NewAtomicLevelAt(zap.DebugLevel),
		"info":  zap.NewAtomicLevelAt(zap.InfoLevel),
		"warn":  zap.NewAtomicLevelAt(zap.WarnLevel),
		"error": zap.NewAtomicLevelAt(zap.ErrorLevel),
		"panic": zap.NewAtomicLevelAt(zap.PanicLevel),
		"fatal": zap.NewAtomicLevelAt(zap.FatalLevel),
	}
	loggingLevel, ok := levels[level]
	if !ok {
		loggingLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}
	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    "func",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(os.Stdout),
		loggingLevel,
	)).Sugar()
}

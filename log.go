package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	LoggerEnv  = "LOGGER_ENV"
	LoggerName = "LOGGER_NAME"
	stdOut     = "stdout"
)

var logger *zap.SugaredLogger

func init() {
	if os.Getenv(LoggerEnv) == "development" {
		logger = newLogger("console")
	} else {
		logger = newLogger("json")
	}
	if os.Getenv(LoggerName) != "" {
		logger = logger.Named(os.Getenv(LoggerName))
	}
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Get() *zap.Logger {
	return logger.Desugar()
}

func newLogger(encoding string) *zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:         encoding,
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{stdOut},
		ErrorOutputPaths: []string{stdOut},
		EncoderConfig: zapcore.EncoderConfig{
			NameKey:        "name",
			MessageKey:     "message",
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			TimeKey:        "time",
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
		},
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return l.WithOptions(
		zap.AddCallerSkip(2),
	).Sugar()
}

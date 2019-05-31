package main

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(debug bool) *zap.Logger {
	encConf := zap.NewDevelopmentEncoderConfig()
	encConf.EncodeDuration = zapcore.StringDurationEncoder
	encConf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04"))
	}

	conf := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "console",
		EncoderConfig:    encConf,
		DisableCaller:    true,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	if debug {
		conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		conf.Development = true
		conf.DisableCaller = false
	}

	logger, err := conf.Build()
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}

	return logger
}

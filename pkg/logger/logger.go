package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/lumberjack.v2"
)

type ContextKey struct{}

type Config struct {
	Level    string
	Env      string
	FilePath string
}

func New(cfg Config) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	encoder := buildEncoder(cfg.Env)
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	}

	if cfg.FilePath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level))
	}

	core := zapcore.NewTee(cores...)

	opts := []zap.Option{zap.AddCaller()}
	if cfg.Env == "dev" {
		opts = append(opts, zap.Development())
	}
	return zap.New(core, opts...), nil
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ContextKey{}).(*zap.Logger); ok {
		return logger
	}
	return zap.L()
}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextKey{}, logger)
}

func buildEncoder(env string) zapcore.Encoder {
	if env == "dev" {
		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(cfg)
	}
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func maxOrDefault(value, defaultValue int) int {
	if value > 0 {
		return value
	}
	return defaultValue
}

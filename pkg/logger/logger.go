package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"isscan/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger(cfg config.LogConfig) error {
	// 创建日志目录
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 设置日志级别
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	var cores []zapcore.Core

	// 普通日志
	normalLog := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Path, cfg.Filename),
		MaxSize:    cfg.MaxSize,
		MaxBackups: 3,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	normalCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(normalLog),
		level,
	)
	cores = append(cores, normalCore)

	// 错误日志
	errorLog := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Path, cfg.ErrorFilename),
		MaxSize:    cfg.MaxSize,
		MaxBackups: 3,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(errorLog),
		zap.ErrorLevel,
	)
	cores = append(cores, errorCore)

	// 控制台输出
	if cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 创建logger
	core := zapcore.NewTee(cores...)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Debug 输出debug级别日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info 输出info级别日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 输出warn级别日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 输出error级别日志
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal 输出fatal级别日志
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// 添加结构化日志
type Field struct {
	Key   string
	Value interface{}
}

func WithTraceID(traceID string) Field {
	return Field{Key: "trace_id", Value: traceID}
}

func WithDuration(d time.Duration) Field {
	return Field{Key: "duration", Value: d.Milliseconds()}
}

// 添加上下文日志
func FromContext(ctx context.Context) *zap.Logger {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return logger.With(zap.String("trace_id", traceID))
	}
	return logger
}

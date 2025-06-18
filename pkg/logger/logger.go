package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log 全局日志实例
	Log *logrus.Logger
)

// Config 日志配置
type Config struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// InitLogger 初始化日志
func InitLogger(config *Config) error {
	Log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// 设置日志格式
	if config.Format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置输出
	if config.Output != "stdout" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.Output)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// 使用 lumberjack 进行日志轮转
		writer := &lumberjack.Logger{
			Filename:   config.Output,
			MaxSize:    config.MaxSize,    // MB
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,     // days
			Compress:   config.Compress,
		}
		Log.SetOutput(writer)
	}

	return nil
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	return Log
}

// Debug 调试日志
func Debug(args ...interface{}) {
	if Log != nil {
		Log.Debug(args...)
	}
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	if Log != nil {
		Log.Debugf(format, args...)
	}
}

// Info 信息日志
func Info(args ...interface{}) {
	if Log != nil {
		Log.Info(args...)
	}
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	if Log != nil {
		Log.Infof(format, args...)
	}
}

// Warn 警告日志
func Warn(args ...interface{}) {
	if Log != nil {
		Log.Warn(args...)
	}
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	if Log != nil {
		Log.Warnf(format, args...)
	}
}

// Error 错误日志
func Error(args ...interface{}) {
	if Log != nil {
		Log.Error(args...)
	}
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	if Log != nil {
		Log.Errorf(format, args...)
	}
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	if Log != nil {
		Log.Fatal(args...)
	}
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	if Log != nil {
		Log.Fatalf(format, args...)
	}
}

// WithField 添加字段
func WithField(key string, value interface{}) *logrus.Entry {
	if Log != nil {
		return Log.WithField(key, value)
	}
	return nil
}

// WithFields 添加多个字段
func WithFields(fields logrus.Fields) *logrus.Entry {
	if Log != nil {
		return Log.WithFields(fields)
	}
	return nil
}

// WithError 添加错误
func WithError(err error) *logrus.Entry {
	if Log != nil {
		return Log.WithError(err)
	}
	return nil
} 
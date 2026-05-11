package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

// InitLogger 初始化日志
// level: 日志级别 (debug, info, warn, error)
// output: 输出目标 (stdout, stderr, 或文件路径)
func InitLogger(level string, output string) {
	var writer io.Writer

	switch output {
	case "stdout":
		// 屏幕输出使用 console 格式化（彩色）
		writer = zerolog.ConsoleWriter{
			Out:        colorable.NewColorableStdout(),
			TimeFormat: time.RFC3339,
		}
	case "stderr":
		writer = zerolog.ConsoleWriter{
			Out:        colorable.NewColorableStderr(),
			TimeFormat: time.RFC3339,
		}
	default:
		// 文件输出保持 JSON 格式
		file, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			writer = zerolog.ConsoleWriter{
				Out:        colorable.NewColorableStdout(),
				TimeFormat: time.RFC3339,
			}
		} else {
			writer = file
		}
	}

	initLoggerWithWriter(level, writer)
}

// InitLoggerMulti 初始化日志，同时输出到多个目标
// level: 日志级别
// outputs: 输出目标列表，支持 "stdout", "stderr", 或文件路径
func InitLoggerMulti(level string, outputs []string) {
	var writers []io.Writer

	for _, output := range outputs {
		var writer io.Writer
		switch output {
		case "stdout":
			// 屏幕输出使用 console 格式化（彩色）
			writer = zerolog.ConsoleWriter{
				Out:        colorable.NewColorableStdout(),
				TimeFormat: time.RFC3339,
			}
		case "stderr":
			writer = zerolog.ConsoleWriter{
				Out:        colorable.NewColorableStderr(),
				TimeFormat: time.RFC3339,
			}
		default:
			// 确保目录存在
			dir := filepath.Dir(output)
			if dir != "." && dir != "" {
				os.MkdirAll(dir, 0755)
			}
			// 文件输出保持 JSON 格式
			file, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				writer = zerolog.ConsoleWriter{
					Out:        colorable.NewColorableStdout(),
					TimeFormat: time.RFC3339,
				}
			} else {
				writer = file
			}
		}
		writers = append(writers, writer)
	}

	// 使用 MultiLevelWriter 合并多个输出
	multi := zerolog.MultiLevelWriter(writers...)
	initLoggerWithWriter(level, multi)
}

// InitDefaultLogger 使用默认配置初始化日志
// 默认级别: info，默认输出: stdout
func InitDefaultLogger() {
	InitLogger("info", "stdout")
}

// InitDefaultLoggerMulti 使用默认配置初始化日志，同时输出到屏幕和文件
// 默认级别: info，输出: stdout + log/app.log
func InitDefaultLoggerMulti() {
	InitLoggerMulti("info", []string{"stdout", "log/app.log"})
}

func initLoggerWithWriter(level string, writer io.Writer) {
	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	logger = zerolog.New(writer).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Debug 调试日志（格式化字符串）
func Debug(msg string, args ...interface{}) {
	logger.Debug().Msgf(msg, args...)
}

// Info 信息日志（格式化字符串）
func Info(msg string, args ...interface{}) {
	logger.Info().Msgf(msg, args...)
}

// Warn 警告日志（格式化字符串）
func Warn(msg string, args ...interface{}) {
	logger.Warn().Msgf(msg, args...)
}

// Error 错误日志（格式化字符串）
func Error(msg string, args ...interface{}) {
	logger.Error().Msgf(msg, args...)
}

// Fatal 致命错误日志（格式化字符串，会调用 os.Exit(1)）
func Fatal(msg string, args ...interface{}) {
	logger.Fatal().Msgf(msg, args...)
}

// Panic 恐慌日志（格式化字符串，会 panic）
func Panic(msg string, args ...interface{}) {
	logger.Panic().Msgf(msg, args...)
}

// DebugWithFields 调试日志（带字段）
func DebugWithFields(msg string, fields ...interface{}) {
	logger.Debug().Fields(toFields(fields)).Msg(msg)
}

// InfoWithFields 信息日志（带字段）
func InfoWithFields(msg string, fields ...interface{}) {
	logger.Info().Fields(toFields(fields)).Msg(msg)
}

// WarnWithFields 警告日志（带字段）
func WarnWithFields(msg string, fields ...interface{}) {
	logger.Warn().Fields(toFields(fields)).Msg(msg)
}

// ErrorWithFields 错误日志（带字段）
func ErrorWithFields(msg string, fields ...interface{}) {
	logger.Error().Fields(toFields(fields)).Msg(msg)
}

// Debugf 格式化调试日志（同 Debug）
func Debugf(format string, args ...interface{}) {
	logger.Debug().Msgf(format, args...)
}

// Infof 格式化信息日志（同 Info）
func Infof(format string, args ...interface{}) {
	logger.Info().Msgf(format, args...)
}

// Warnf 格式化警告日志（同 Warn）
func Warnf(format string, args ...interface{}) {
	logger.Warn().Msgf(format, args...)
}

// Errorf 格式化错误日志（同 Error）
func Errorf(format string, args ...interface{}) {
	logger.Error().Msgf(format, args...)
}

// Fatalln 致命错误日志（带换行）
func Fatalln(msg string) {
	logger.Fatal().Msg(msg)
}

// Print 打印日志（相当于 Info）
func Print(msg string) {
	logger.Info().Msg(msg)
}

// Printf 格式化打印日志
func Printf(format string, args ...interface{}) {
	logger.Info().Msgf(format, args...)
}

// WithError 记录错误
func WithError(err error) *zerolog.Event {
	return logger.Error().Err(err)
}

// WithField 添加单个字段
func WithField(key string, value interface{}) *zerolog.Event {
	return logger.Info().Fields(map[string]interface{}{key: value})
}

// WithFields 添加多个字段
func WithFields(fields map[string]interface{}) *zerolog.Event {
	return logger.Info().Fields(fields)
}

// toFields 将可变参数转换为 map
func toFields(fields []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				result[key] = fields[i+1]
			}
		}
	}
	return result
}

/*
使用示例：

package main

import (
    "server/pkg/log"
)

func main() {
    // 初始化日志：同时输出到屏幕（彩色）和文件（JSON）
    log.InitLoggerMulti("debug", []string{"stdout", "log/app.log"})

    // 格式化字符串方式（默认）
    log.Info("Server started on port %d", 8080)
    log.Debug("User %s logged in", username)
    log.Warn("Low memory warning: %d MB", 100)
    log.Error("Failed to connect: %v", err)

    // 带字段方式
    log.InfoWithFields("User login", "user_id", 123, "username", "alice")
    log.ErrorWithFields("Database error", "error", err, "query", query)
}

输出示例：
屏幕输出（彩色）:
  2026-05-11T15:35:18+08:00 INF Server started on port 8080

文件输出（JSON）:
  {"level":"info","time":"2026-05-11T15:35:18+08:00","caller":"main.go:42","message":"Server started on port 8080"}
*/

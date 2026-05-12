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

func InitLogger(level string, output string) {
	var writer io.Writer

	switch output {
	case "stdout":
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

func InitLoggerMulti(level string, outputs []string) {
	var writers []io.Writer

	for _, output := range outputs {
		var writer io.Writer
		switch output {
		case "stdout":
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
			dir := filepath.Dir(output)
			if dir != "." && dir != "" {
				os.MkdirAll(dir, 0755)
			}
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

	multi := zerolog.MultiLevelWriter(writers...)
	initLoggerWithWriter(level, multi)
}

func InitDefaultLogger() {
	InitLogger("info", "stdout")
}

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
		Logger()
}

func Debug(msg string, fields ...interface{}) {
	logger.Debug().Caller(1).Fields(toFields(fields)).Msg(msg)
}

func Info(msg string, fields ...interface{}) {
	logger.Info().Caller(1).Fields(toFields(fields)).Msg(msg)
}

func Warn(msg string, fields ...interface{}) {
	logger.Warn().Caller(1).Fields(toFields(fields)).Msg(msg)
}

func Error(msg string, fields ...interface{}) {
	logger.Error().Caller(1).Fields(toFields(fields)).Msg(msg)
}

func Fatal(msg string, fields ...interface{}) {
	logger.Fatal().Caller(1).Fields(toFields(fields)).Msg(msg)
}

func Panic(msg string, fields ...interface{}) {
	logger.Panic().Caller(1).Fields(toFields(fields)).Msg(msg)
}

func Debugf(format string, args ...interface{}) {
	logger.Debug().Caller(1).Msgf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Info().Caller(1).Msgf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warn().Caller(1).Msgf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Error().Caller(1).Msgf(format, args...)
}

func Fatalln(msg string) {
	logger.Fatal().Caller(1).Msg(msg)
}

func Print(msg string) {
	logger.Info().Caller(1).Msg(msg)
}

func Printf(format string, args ...interface{}) {
	logger.Info().Caller(1).Msgf(format, args...)
}

func WithError(err error) *zerolog.Event {
	return logger.Error().Err(err)
}

func WithField(key string, value interface{}) *zerolog.Event {
	return logger.Info().Fields(map[string]interface{}{key: value})
}

func WithFields(fields map[string]interface{}) *zerolog.Event {
	return logger.Info().Fields(fields)
}

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

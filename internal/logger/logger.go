package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/vviveksharma/lorebound/internal/utils"
)

var (
	once sync.Once
	log  *slog.Logger
)

// Init configures the process-wide default slog logger from LOG_LEVEL.
func Init() {
	once.Do(func() {
		level := parseLevel(utils.Getenv("LOG_LEVEL", "info"))
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
		slog.SetDefault(log)
	})
}

func ensure() {
	if log == nil {
		Init()
	}
}

// L returns the shared logger instance.
func L() *slog.Logger {
	ensure()
	return log
}

// Writer is an io.Writer that forwards lines to Info (e.g. HTTP access logs).
func Writer() io.Writer {
	return logWriter{}
}

type logWriter struct{}

func (logWriter) Write(p []byte) (int, error) {
	ensure()
	line := strings.TrimSpace(string(p))
	if line != "" {
		log.Info(line)
	}
	return len(p), nil
}

func Info(msg string, args ...any) {
	ensure()
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	ensure()
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	ensure()
	log.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	ensure()
	log.Debug(msg, args...)
}

func Fatal(msg string, args ...any) {
	ensure()
	log.Error(msg, args...)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	ensure()
	log.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

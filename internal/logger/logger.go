package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// AppLogger — реализация Logger с использованием logrus.
type AppLogger struct {
	*logrus.Logger
}

// Реализация методов интерфейса logger.Logger
func (l *AppLogger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *AppLogger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *AppLogger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *AppLogger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *AppLogger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

func (l *AppLogger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

func (l *AppLogger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

func (l *AppLogger) SetLevel(level logrus.Level) {
	l.Logger.SetLevel(level)
}

// NewLogger создаёт новый экземпляр логгера.
func NewLogger(logFile string) *AppLogger {
	logger := logrus.New()
	// Формат: JSON для удобства парсинга
	logger.SetFormatter(&logrus.JSONFormatter{})

	dir := filepath.Dir(logFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	logger.SetOutput(file)

	level := os.Getenv("LOG_LEVEL")
	parsedLevel, _ := logrus.ParseLevel(level)
	logger.SetLevel(parsedLevel)

	return &AppLogger{logger}
}

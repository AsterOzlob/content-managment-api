package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger — интерфейс для унификации логгера в приложении.
type Logger interface {
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields map[string]interface{}) *logrus.Entry
	WithError(err error) *logrus.Entry
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	SetLevel(level logrus.Level)
}

package logging

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Logger представляет настраиваемый логгер.
type Logger struct {
	*logrus.Logger
}

// NewLogger создает новый экземпляр логгера.
func NewLogger(logFile string) *Logger {
	logger := logrus.New()

	// Настройка формата логов
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Получаем директорию из пути к файлу
	dir := filepath.Dir(logFile)

	// Проверяем, существует ли директория, и создаем её, если её нет
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// Создание файла для логов
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	// Направление вывода логов в файл
	logger.SetOutput(file)

	// Уровень логирования
	logger.SetLevel(logrus.DebugLevel)

	return &Logger{logger}
}

// Log записывает сообщение с указанным уровнем логирования и дополнительными полями.
func (l *Logger) Log(level logrus.Level, message string, fields map[string]interface{}) {
	entry := l.WithFields(fields)
	switch level {
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.DebugLevel:
		entry.Debug(message)
	default:
		entry.Info(message) // По умолчанию используем уровень Info
	}
}

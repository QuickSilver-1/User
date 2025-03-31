package logger

import (
	"go.uber.org/zap"
)

// Logger - логгирования с использованием zap
var (
	Logger *zap.Logger
)

// NewLogger создает новый экземпляр логгера с заданным путем для вывода логов
func NewLogger() error {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}

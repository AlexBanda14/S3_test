package logger

import (
	"fmt"
	"go.uber.org/zap"
)

type ZapLogger struct {
	Logger *zap.SugaredLogger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("error creating zap logger: %w", err)
	}
	return &ZapLogger{Logger: logger.Sugar()}, nil
}

func (z *ZapLogger) Close() error {
	err := z.Logger.Sync()
	if err != nil {
		return fmt.Errorf("error closing zap logger: %w", err)
	}
	return nil
}

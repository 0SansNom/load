package logger_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"yourmodule/internal/logger"
)

func TestNewLogger_Production(t *testing.T) {
	log, err := logger.New(true)
	require.NoError(t, err)
	require.IsType(t, &zap.Logger{}, log)
}

func TestNewLogger_Development(t *testing.T) {
	log, err := logger.New(false)
	require.NoError(t, err)
	require.IsType(t, &zap.Logger{}, log)
}

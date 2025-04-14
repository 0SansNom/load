package logger

import (
	"go.uber.org/zap"
)

// New returns a new zap.Logger instance.
// Use production mode by default; development mode can be toggled if needed.
func New(production bool) (*zap.Logger, error) {
	if production {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

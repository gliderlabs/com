package zap

import (
	"github.com/gliderlabs/com"
	"go.uber.org/zap"
)

// TODO: wrap for With
// TODO: component for config

// Register the zap logger component with a registry
func Register(registry *com.Registry) error {
	logger, _ := zap.NewDevelopment()
	return registry.Register(&com.Object{Value: logger.Sugar()})
}

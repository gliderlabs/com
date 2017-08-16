package zap

import (
	"github.com/gliderlabs/com"
	"go.uber.org/zap"
)

func Register(registry *com.Registry) {
	logger, _ := zap.NewDevelopment()
	registry.Register(&com.Object{Value: logger.Sugar()})
}

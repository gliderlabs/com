package foobar

import (
	"time"

	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/com/log"
	"github.com/progrium/cmd/com/web"
)

type Config struct {
	Foobar  string
	Timeout time.Duration
}

type Component struct {
	Config
	Log     log.DebugLogger `com:"singleton"`
	Web     *web.Component  `com:"singleton"`
	Handler Handler         `com:"config"`
}

type Handler interface {
	HandleFoobar() string
}

func (com *Component) InitializeDaemon() error {
	com.Log.Debugw("foobar init", "from", "foobar")
	return nil
}

func (com *Component) Serve() {
	for {
		com.Log.Debugf(com.Handler.HandleFoobar())
		time.Sleep(5 * time.Second)
	}
}

func (com *Component) Stop() {
	com.Log.Debugf("foobar stop")
}

func (com *Component) InitializeConfig(config config.Settings) error {
	config.SetDefault("foobar", "value")
	config.SetDefault("timeout", 10*time.Second)
	config.SetDefault("handler", "qux")
	return config.Load(&(com.Config))
}

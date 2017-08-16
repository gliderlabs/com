package qux

import (
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/example/app/foobar"
	"github.com/gliderlabs/com/log"
)

func init() {
	err := com.DefaultRegistry.Register(
		&com.Object{Value: &Component{}})
	if err != nil {
		panic(err)
	}
}

type Component struct {
	Log    log.DebugLogger   `com:"singleton"`
	Foobar *foobar.Component `com:"singleton"`
}

func (com *Component) InitializeDaemon() error {
	com.Log.Debugf("baz init")
	return nil
}

func (com *Component) HandleFoobar() string {
	return "baz handler"
}

func (com *Component) Serve() {
	for {
	}
}

func (com *Component) Stop() {
	com.Log.Debugf("qux stop")
}

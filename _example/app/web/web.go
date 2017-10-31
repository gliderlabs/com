package web

import (
	"net/http"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/com/log"
)

func init() {
	com.DefaultRegistry.Register(&com.Object{Value: &Component{}})
}

type Config struct {
	Addr string `mapstructure:"listen_addr"`
}

type Component struct {
	Log log.Logger `com:"singleton"`

	server *http.Server

	Config
}

func (c *Component) InitializeConfig(config config.Settings) error {
	return config.Unmarshal(&(c.Config))
}

func (c *Component) InitializeDaemon() error {
	c.server = &http.Server{
		Addr:    c.Addr,
		Handler: c,
	}
	return nil
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world\n"))
}

func (c *Component) Serve() {
	c.Log.Info("serving on", c.Addr, "...")
	if err := c.server.ListenAndServe(); err != nil {
		c.Log.Info(err)
	}
}

func (c *Component) Stop() {
	c.Log.Info("stopping...")
}

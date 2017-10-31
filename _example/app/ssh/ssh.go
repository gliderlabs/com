package ssh

import (
	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/com/log"
	"github.com/gliderlabs/ssh"
)

type Config struct {
	Addr string `mapstructure:"listen_addr"`
}

type Component struct {
	Log log.Logger `com:"singleton"`

	server *ssh.Server
	Config
}

func (c *Component) InitializeConfig(config config.Settings) error {
	return config.Unmarshal(&(c.Config))
}

func (c *Component) InitializeDaemon() error {
	c.server = &ssh.Server{
		Addr:    c.Addr,
		Handler: c.HandleSession,
	}
	return nil
}

func (c *Component) HandleSession(sess ssh.Session) {
	sess.Write([]byte("Hello world\n"))
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

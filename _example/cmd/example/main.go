package main

import (
	"log"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/com/config/viper"
	"github.com/gliderlabs/com/daemon"
	stdlog "github.com/gliderlabs/com/log/std"
	"github.com/gliderlabs/com/plugins"
)

const name = "example"

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	assert(stdlog.Register(com.DefaultRegistry))
	assert(plugins.Load(com.DefaultRegistry, name, []string{"."}))
	assert(config.Load(com.DefaultRegistry, viper.New(), name, []string{"."}))
	assert(daemon.Run(com.DefaultRegistry, name))
}

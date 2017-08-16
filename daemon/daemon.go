package daemon

import (
	"log"
	"os"
	"os/signal"
	"path"
	"strings"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config/viper"
	"github.com/gliderlabs/com/log/std"
	"github.com/thejerf/suture"
)

type Component struct {
	Initializers []Initializer    `com:"extpoint"`
	Services     []suture.Service `com:"extpoint"`
	Terminators  []Terminator     `com:"extpoint"`
}

type Initializer interface {
	InitializeDaemon() error
}

type Terminator interface {
	TerminateDaemon() error
}

func Run(name string) {
	daemon := &Component{}
	std.Register(com.DefaultRegistry)

	// initial registry population
	if err := com.DefaultRegistry.Register(&com.Object{Value: daemon}); err != nil {
		log.Fatal(err)
	}

	// setup configuration
	configName := name
	configPath := "."
	configFilepath := os.Getenv(strings.ToUpper(name) + "_CONFIG")
	if configFilepath != "" {
		configPath = path.Dir(configFilepath)
		configName = strings.TrimSuffix(path.Base(configFilepath), path.Ext(configFilepath))
	}
	if err := viper.Load(com.DefaultRegistry, configName, []string{configPath}); err != nil {
		log.Fatal(err)
	}

	// call initializers
	for _, i := range daemon.Initializers {
		if err := i.InitializeDaemon(); err != nil {
			log.Fatal(err)
		}
	}

	// setup all services
	app := suture.NewSimple(name)
	for _, s := range daemon.Services {
		app.Add(s)
	}

	// setup terminators on SIGINT
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	go func() {
		<-sigInt
		app.Stop()
		for _, i := range daemon.Terminators {
			if err := i.TerminateDaemon(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// serve all
	app.Serve()
}

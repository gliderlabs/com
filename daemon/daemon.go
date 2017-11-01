package daemon

import (
	"errors"
	"os"
	"os/signal"

	"github.com/gliderlabs/com"
	"github.com/thejerf/suture"
)

var (
	// ErrNoServices is returned when there are no objects in the registry to start as services
	ErrNoServices = errors.New("no services in registry to run")
)

// Component is the main package component.
type Component struct {
	Initializers []Initializer    `com:"extpoint"`
	Services     []suture.Service `com:"extpoint"`
	Terminators  []Terminator     `com:"extpoint"`
}

// Initializer is an extension point interface with a hook for components to
// initialize before services are started. Returning an error will cancel
// the start of the daemon services.
type Initializer interface {
	InitializeDaemon() error
}

// Terminator is an extension point interface with a hook for components to
// handle daemon termination via SIGINT signals, which will first Stop all
// services and then run this hook.
type Terminator interface {
	TerminateDaemon() error
}

// Run takes a registry and registers the daemon.Component, calls Initializers,
// sets up a parent suture.Service adding any registered suture.Services to it,
// sets up a goroutine to listen for SIGINT, then calls the blocking Serve on
// the parent suture.Service. The first error returned by an Initializer or
// Terminator is returned if any.
func Run(registry *com.Registry, name string) error {
	daemon := &Component{}

	// initial registry population
	if err := registry.Register(&com.Object{Value: daemon}); err != nil {
		return err
	}

	// call initializers
	for _, i := range daemon.Initializers {
		if err := i.InitializeDaemon(); err != nil {
			return err
		}
	}

	// fail if no services
	if len(daemon.Services) == 0 {
		return ErrNoServices
	}

	// setup all services
	app := suture.NewSimple(name)
	for _, s := range daemon.Services {
		app.Add(s)
	}

	// setup terminators on SIGINT
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	termErr := make(chan error)
	go func() {
		<-sigInt
		app.Stop()
		for _, i := range daemon.Terminators {
			if err := i.TerminateDaemon(); err != nil {
				termErr <- err
				break
			}
		}
		close(termErr)
	}()

	// serve all
	app.Serve()

	return <-termErr
}

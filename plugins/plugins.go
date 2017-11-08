package plugins

import (
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/gliderlabs/com/objects"
	"github.com/spf13/afero"
)

const (
	// TODO: should be public to self document
	envFormatter = "%s_PLUGINS"

	// TODO: should be public to self document
	fnSymbol = "Registerable"
)

var (
	// used for testing. TODO: component?
	pluginLoader = loadPlugin
	fs           = afero.NewOsFs()
)

// TODO: define and use a Registerable type to self document the function signature

// Load will open Go shared object plugins, call the symbol Registerable, and
// register objects returned.
func Load(registry *objects.Registry, name string, paths []string) error {
	// get paths from environment
	envPaths := os.Getenv(fmt.Sprintf(envFormatter, strings.ToUpper(name)))
	paths = append(paths, strings.Split(envPaths, ":")...)

	// look for .so files in each path and try to load them
	for _, path := range paths {
		matches, err := afero.Glob(fs, fmt.Sprintf("%s/*.so", path))
		if err != nil {
			continue
		}
		for _, filepath := range matches {
			if err := pluginLoader(registry, filepath); err != nil {
				return err
			}
		}
	}

	return nil
}

func loadPlugin(reg *objects.Registry, filepath string) error {
	// open the shared object binary
	p, err := plugin.Open(filepath)
	if err != nil {
		return err
	}

	// lookup function that returns objects to register
	symbol, err := p.Lookup(fnSymbol)
	if err != nil {
		return err
	}

	// call function and register returned objects
	for _, obj := range symbol.(func() []interface{})() {
		reg.Register(&objects.Object{Value: obj})
	}

	return nil
}

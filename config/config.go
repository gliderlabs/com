package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/gliderlabs/com/objects"
)

const (
	envFormatter = "%s_CONFIG"
	disabledKey  = "disabled"
)

// Load uses a Provider to read in configuration from various files and the
// environment for the objects in a particular Registry. It passes Settings for
// each object to any object that implements the Initializer interface. Then it
// will lookup any struct fields with `com:"config"` and use the settings for
// that object to get the name of an object from the registry to assign to that
// field. It also disables any objects in the Registry referenced in the top-level
// config section called "disabled".
func Load(registry *objects.Registry, provider Provider, name string, paths []string) error {
	// add extra paths from environment
	envConfig := os.Getenv(fmt.Sprintf(envFormatter, strings.ToUpper(name)))
	paths = append(paths, strings.Split(envConfig, ":")...)

	// tell provider to load config
	cfg, err := provider.Load(name, paths)
	if err != nil {
		return err
	}

	// get all top level keys in config
	var keys map[string]interface{}
	if err := cfg.Unmarshal(&keys); err != nil {
		return err
	}

	// iterate over all objects in registry
	for _, obj := range registry.Objects() {
		s := provider.New()

		// check if any top level key matches object
		for key := range keys {
			o, _ := registry.Lookup(key)
			if o == obj {
				s = cfg.Sub(key)
				break
			}
		}

		// if object is config.Initializer, initialize it
		if init, ok := obj.Value.(Initializer); ok {
			if err := init.InitializeConfig(s); err != nil {
				return err
			}
		}

		// use config to lookup and set config fields
		for name, field := range obj.Fields {
			if field.Config && s.IsSet(name) {
				objName, ok := s.Get(name).(string)
				if !ok {
					continue
				}
				o, err := registry.Lookup(objName)
				if err != nil {
					return err
				}
				obj.Assign(name, o)
			}
		}
	}

	// disable any objects found under disabled key
	if cfg.IsSet(disabledKey) {
		var disabled map[string]bool
		cfg.UnmarshalKey(disabledKey, &disabled)
		for name, d := range disabled {
			o, err := registry.Lookup(name)
			if err != nil {
				continue
			}
			registry.SetEnabled(o.FQN(), !d)
		}
	}

	// reload registry
	return registry.Reload()
}

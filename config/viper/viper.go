package viper

import (
	"strings"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config"
	"github.com/spf13/viper"
)

type settings struct {
	*viper.Viper
}

func (s *settings) Sub(key string) config.Settings {
	return &settings{s.Viper.Sub(key)}
}

func (s *settings) Load(rawVal interface{}) error {
	return s.Viper.Unmarshal(rawVal)
}

func (s *settings) LoadKey(key string, rawVal interface{}) error {
	return s.Viper.UnmarshalKey(key, rawVal)
}

func Load(registry *com.Registry, name string, paths []string) error {
	cfg := viper.New()
	if len(paths) > 0 && name != "" {
		cfg.SetConfigName(name)
		for _, p := range paths {
			cfg.AddConfigPath(p)
		}
		if err := cfg.ReadInConfig(); err != nil {
			return err
		}
	}
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	var keys map[string]interface{}
	if err := cfg.Unmarshal(&keys); err != nil {
		return err
	}
	for _, obj := range registry.Objects() {
		s := &settings{viper.New()}
		for key, _ := range keys {
			o, _ := registry.Lookup(key)
			if o == obj {
				s = &settings{cfg.Sub(key)}
				break
			}
		}
		if init, ok := obj.Value.(config.Initializer); ok {
			if err := init.InitializeConfig(s); err != nil {
				return err
			}
		}
		for name, field := range obj.Fields {
			if field.Config && s.IsSet(name) {
				o, err := registry.Lookup(s.GetString(name))
				if err != nil {
					return err
				}
				obj.Assign(name, o)
			}
		}
	}
	if cfg.IsSet("disabled") {
		var disabled map[string]bool
		cfg.UnmarshalKey("disabled", &disabled)
		for fqn, d := range disabled {
			registry.SetEnabled(fqn, !d)
		}
	}
	return registry.Reload()
}

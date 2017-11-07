package viper

import (
	"strings"

	"github.com/gliderlabs/com/config"
	"github.com/spf13/viper"
)

// New returns an initialized Viper provider instance.
func New() config.Provider {
	v := viper.New()
	return &Provider{v}
}

// Provider is a config.Provider for Viper.
type Provider struct {
	*viper.Viper
}

// Sub returns new Settings instance representing a sub tree of this instance.
func (p *Provider) Sub(key string) config.Settings {
	sub := p.Viper.Sub(key)
	// Sub somehow removes values set/overridden by environment, so we do this
	// just to make sure they are set properly in this new Viper instance.
	// UnmarshalKey could be used here except it has the same problem as Sub.
	var keys map[string]map[string]interface{}
	p.Unmarshal(&keys)
	for k, v := range keys[key] {
		sub.Set(k, v)
	}
	return &Provider{sub}
}

// New returns an empty Settings instance.
func (p *Provider) New() config.Settings {
	return New()
}

// Load returns Settings for named configuration loaded from provided paths.
// Leave the file extension off of name as supported format extensions will
// automatically be added by Viper.
func (p *Provider) Load(name string, paths []string) (config.Settings, error) {
	// read in config files
	if len(paths) > 0 && name != "" {
		p.SetConfigName(name)
		for _, path := range paths {
			p.AddConfigPath(path)
		}
		if err := p.ReadInConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				// it's fine!
				break
			default:
				return nil, err
			}
		}
	}

	// read config from environment
	p.AutomaticEnv()
	p.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return p, nil
}

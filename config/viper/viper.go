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

type Provider struct {
	*viper.Viper
}

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

func (p *Provider) Empty() config.Settings {
	return New()
}

func (p *Provider) Load(name string, paths []string) (config.Settings, error) {
	// read in config files
	if len(paths) > 0 && name != "" {
		p.SetConfigName(name)
		for _, path := range paths {
			p.AddConfigPath(path)
		}
		if err := p.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	// read config from environment
	p.AutomaticEnv()
	p.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return p, nil
}

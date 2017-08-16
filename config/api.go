package config

import "time"

type Settings interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	IsSet(key string) bool

	Load(rawVal interface{}) error
	LoadKey(key string, rawVal interface{}) error

	SetDefault(key string, value interface{})
	Sub(key string) Settings
}

type Initializer interface {
	InitializeConfig(config Settings) error
}

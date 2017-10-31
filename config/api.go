package config

import "time"

// Settings is an interface representing a collection of key-values from a
// configuration or subset of a configuration.
type Settings interface {
	// Get returns the value associated with the key as an empty interface.
	Get(key string) interface{}

	// GetBool returns the value associated with the key as a boolean.
	GetBool(key string) bool

	// GetFloat64 returns the value associated with the key as a float64.
	GetFloat64(key string) float64

	// GetInt returns the value associated with the key as an integer.
	GetInt(key string) int

	// GetString returns the value associated with the key as a string.
	GetString(key string) string

	// GetStringMap returns the value associated with the key as a map of interfaces.
	GetStringMap(key string) map[string]interface{}

	// GetStringMapString returns the value associated with the key as a map of strings.
	GetStringMapString(key string) map[string]string

	// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
	GetStringMapStringSlice(key string) map[string][]string

	// GetStringSlice returns the value associated with the key as a slice of strings.
	GetStringSlice(key string) []string

	// GetTime returns the value associated with the key as time.
	GetTime(key string) time.Time

	// GetDuration returns the value associated with the key as a duration.
	GetDuration(key string) time.Duration

	// IsSet checks to see if the key has been set in configuration.
	IsSet(key string) bool

	// Unmarshal unmarshals the config into a Struct. Make sure that the tags on the fields of the structure are properly set.
	Unmarshal(rawVal interface{}) error

	// UnmarshalKey takes a single key and unmarshals it into a Struct.
	// BUG: Currently with Viper, UnmarshalKey will ignore environment values.
	UnmarshalKey(key string, rawVal interface{}) error

	// SetDefault sets the default value for this key.
	SetDefault(key string, value interface{})

	// Sub returns new Settings instance representing a sub tree of this instance.
	Sub(key string) Settings
}

// Provider is an interface for configuration providers, which includes the
// the interface to loaded configuration from that provider.
type Provider interface {
	// Load returns Settings for named configuration loaded from provided paths.
	// Leave the file extension off of name as supported format extensions will
	// automatically be added.
	Load(name string, paths []string) (Settings, error)

	// Empty returns an empty/new Settings instance.
	Empty() Settings

	Settings
}

// Initializer is an extension point interface with a hook allowing objects to
// handle their configuration when configuration is loaded by the provider.
type Initializer interface {
	// InitializeConfig is called on a registered object with Settings for that
	// object when configuration has been loaded.
	InitializeConfig(config Settings) error
}

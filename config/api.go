package config

// Settings is an interface representing a collection of key-values from a
// configuration or subset of a configuration. 90% of the time you'll just
// use Unmarshal into a struct, but sometimes you'll want to grab a specific
// key. To encourage using structs, there are no typed getters.
type Settings interface {
	// Get returns the value associated with the key as an empty interface.
	Get(key string) interface{}

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

	// New returns an empty Settings instance.
	New() Settings

	Settings
}

// Initializer is an extension point interface with a hook allowing objects to
// handle their configuration when configuration is loaded by the provider.
type Initializer interface {
	// InitializeConfig is called on a registered object with Settings for that
	// object when configuration has been loaded.
	InitializeConfig(config Settings) error
}

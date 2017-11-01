// Package config does 3 things:
//   1. provide an interface for component objects to get and process their own config
//   2. provide an interface for config providers to hook into the load mechanism
//   3. provide a load mechanism with sensible defaults and config lifecycle
//
// The load mechanism defines the basic semantics and lifecycle of configuration:
//
//   1. configuration is loaded from one or more files via Load
//   2. the config format(s) are up to the config provider
//   3. top level keys map to registered object names matched via Lookup
//   4. a special "disabled" top level key is used to disable registered objects
//   5. files are loaded from paths that the app specifies in call to Load
//   6. more filepaths can be specified via user environment variable
//   7. config can be set or overridden by user environment variables
//   8. resulting config for each object is passed via extension point
//   9. objects use this to specify defaults, process, and store values
//   10. "config" fields of an object are assigned by lookup using the key by that field name
//   11. registry is reloaded
//
// The default, preferred, and builtin configuration provider is Viper. Viper
// can be used directly for more control, or replaced with a custom provider.
//
// The Settings and Initializer interfaces are the only parts needed for object
// compatibility in the component ecosystem. Apps can define their own config
// Provider, or ignore the Load mechanism entirely.
package config

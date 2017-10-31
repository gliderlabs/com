// Package viper provides a config.Provider based on the Viper configuration
// library. Since the config.Provider interface was borrowed from Viper, this is
// a very light package implementation.
//
// The Load implementation not only loads config files from multiple paths, it
// uses Viper's AutomaticEnv to load config from environment. It also uses
// SetEnvKeyReplacer to use underscores in place of periods when identifying
// sub keys via environment.
package viper

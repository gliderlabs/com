// Package com is a user-facing interface to the object registry. Since
// the only part of the API you need to use is Register, the rest of the API
// for interacting with an object registry is in its own subpackage. This is used
// by other tooling built around com, for example the config subpackage.
//
// When you register an object, it will populate fields based on the com struct
// tags used on them. The object will then also be used to populate fields
// of other objects in the registry where the type or interface matches.
//
//  type Component struct {
//  	Log      log.Logger     `com:"singleton"`
//  	Handlers []api.Handlers `com:"extpoint"`
//  	DB       api.Store      `com:"config"`
//  }
//
// In the above example component, it has fields with all three possible struct
// tags:
//
// Singleton will pick the first object in the registry that implements that
// interface. You can also use pointers to concrete types, for example to other
// component types.
//
// Extpoint is going to be a slice of all objects in the registry that implement
// that interface.
//
// Config is not populated, but is allowed to be populated via the registry API.
// If you're using the config package, it will do this for you and populate it
// based on configuration. In this case, the key would be "DB" and the value
// could be the name of any registered component that implements api.Store.
package com

import "github.com/gliderlabs/com/objects"

var (
	// DefaultRegistry is often used as the single top level registry for an app.
	DefaultRegistry = &objects.Registry{}
)

// Register will add an object with optional name to the default registry.
func Register(obj interface{}, name string) error {
	return DefaultRegistry.Register(objects.New(obj, name))
}

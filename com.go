package com

import "github.com/gliderlabs/com/registry"

var (
	// DefaultRegistry is often used as the single top level registry for an app
	DefaultRegistry = &registry.Registry{}
)

func Register(obj interface{}, name string) error {
	return DefaultRegistry.Register(&registry.Object{Value: obj, Name: name})
}

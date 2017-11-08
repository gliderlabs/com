package com

import "github.com/gliderlabs/com/objects"

var (
	// DefaultRegistry is often used as the single top level registry for an app
	DefaultRegistry = &objects.Registry{}
)

func Register(obj interface{}, name string) error {
	return DefaultRegistry.Register(objects.New(obj, name))
}

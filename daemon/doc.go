// Package daemon is a simple service supervisor, currently built around
// the service interface of the suture project: https://github.com/thejerf/suture
//
// The top level function Run takes a registry and orchestrates the running of
// a daemon based on suture.Services implemented in the registry. It also
// has hooks for custom initialize and termination behavior.
//
// Daemon is also a great example of a simple component acting as a
// micro-framework. It's built around an interface representing the idea of
// "services". If daemon is not able to satisfy requirements with the use of
// extension points, you can write your own alternative. As long as it uses the
// same interface and semantics for services, it would be compatible with
// components exposing services for this package.
package daemon

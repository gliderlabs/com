# gliderlabs/com

A component-oriented approach to building Golang applications

[![GoDoc](https://godoc.org/github.com/gliderlabs/com?status.svg)](https://godoc.org/github.com/gliderlabs/com)
[![CircleCI](https://img.shields.io/circleci/project/github/gliderlabs/com.svg)](https://circleci.com/gh/gliderlabs/com)*
[![Go Report Card](https://goreportcard.com/badge/github.com/gliderlabs/com)](https://goreportcard.com/report/github.com/gliderlabs/com)
[![Slack](http://slack.gliderlabs.com/badge.svg)](http://slack.gliderlabs.com)
[![Email Updates](https://img.shields.io/badge/updates-subscribe-yellow.svg)](https://app.convertkit.com/landing_pages/289455)

_* Build is failing because of an [upstream issue](https://github.com/gliderlabs/com/issues/1) with a [PR](https://github.com/spf13/viper/pull/405) waiting to be merged_

We want to see a world with great "building blocks" where you can quickly build
whatever you want. Simple and composable is not enough, they need to integrate
and hook into each other.

This library provides the core mechanisms needed to build out a component
architecture for your applications that also extend into an ecosystem of reusable
components.

There are two parts to this package that are designed to work with each other:

 * An object registry for interface-based extension points and dependency injection
 * A configuration API for settings, disabling objects, and picking interface backends

The API and even core functionality alone doesn't imply how we build components
with this tool. We can formalize much of this in a proper framework project once
we've determined the conventions. Until then, we can focus on examples and our
small but growing library of [standard components](https://github.com/gliderlabs/stdcom).

In the end, this package helps facilitates structuring Go programs into modular and
extensible components that become much more drop-in building blocks than
the usual Go package.

More documentation soon.

## Dependencies

Good libraries should have minimal dependencies. Here are the ones com uses and
for what:

 * github.com/spf13/afero (plugins, config tests)
 * github.com/spf13/viper (config, config/viper)

## License

BSD

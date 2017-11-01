# gliderlabs/com

Component framework for Go

[![GoDoc](https://godoc.org/github.com/gliderlabs/com?status.svg)](https://godoc.org/github.com/gliderlabs/com)
[![CircleCI](https://img.shields.io/circleci/project/github/gliderlabs/com.svg)](https://circleci.com/gh/gliderlabs/com)
[![Go Report Card](https://goreportcard.com/badge/github.com/gliderlabs/com)](https://goreportcard.com/report/github.com/gliderlabs/com)
[![Slack](http://slack.gliderlabs.com/badge.svg)](http://slack.gliderlabs.com)
[![Email Updates](https://img.shields.io/badge/updates-subscribe-yellow.svg)](https://app.convertkit.com/landing_pages/289455)

This package helps you organize your Go programs into logical components in a way
that improves:

 * Testability
 * Extensibility
 * Configurability
 * Reuseability

More information soon.


## Dependencies

Good libraries should have minimal dependencies. Here are the ones com uses and
for what:

 * github.com/spf13/afero (plugins, config tests)
 * github.com/spf13/viper (config, config/viper)
 * github.com/thejerf/suture (daemon)
 * go.uber.org/zap (log/zap)

## License

BSD

// Package log defines a consistent logging API that can be used from registered
// objects. This gives the component ecosystem the option to perform logging
// without deciding on the logger that must be used. It also allows apps to
// switch out loggers in different scenarios.
//
// There are several wrapper packages builtin. Since the API is extracted from
// the Zap logger, a convenience Zap wrapper is provided. A wrapper for the
// standard log package is provided, which adds the key-value API Zap introduces.
// There is also a null implementation for testing and development.
//
// There are three main opinions this API comes with:
//   1. Many apps will use and normalize to key-value structured logging.
//   2. Logging shouldn't impact control, so Fatal et al are not used.
//   3. Warning is a useless log level. Ideally only Debug and Info are used.
//
// Using only Debug and Info is highly encouraged, but Error is also provided
// as a compromise in supporting the traditional log level model. Open source
// components should avoid using it. In fact, most open source libraries should
// just use Debug to minimize forced user-facing logs.
//
// This leaves the semantics of Error up to the app. If used, we suggest
// reserving it for logging unhandled error objects.
package log

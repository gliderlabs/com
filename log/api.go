package log

// Logger is a full debug, info, and error logger.
type Logger interface {
	With(args ...interface{}) Logger
	DebugLogger
	InfoLogger
	ErrorLogger
}

// DebugLogger is for developer/operator debug logs.
// Libraries should typically only use DebugLogger.
type DebugLogger interface {
	// Debug logs with fmt.Sprint-like arguments.
	Debug(args ...interface{})

	// Debugf logs with fmt.Sprintf-like arguments.
	Debugf(template string, args ...interface{})

	// Debugw logs with additional context using variadic key-value pairs.
	Debugw(msg string, keysAndValues ...interface{})
}

// InfoLogger is for typical user/operator logs.
// Libraries should avoid using InfoLogger.
type InfoLogger interface {
	// Info logs with fmt.Sprint-like arguments.
	Info(args ...interface{})

	// Infof logs with fmt.Sprintf-like arguments.
	Infof(template string, args ...interface{})

	// Infow logs with additional context using variadic key-value pairs.
	Infow(msg string, keysAndValues ...interface{})
}

// ErrorLogger is specifically for errors.
// Libraries should typically never use ErrorLogger.
type ErrorLogger interface {
	// Error logs with fmt.Sprint-like arguments.
	Error(args ...interface{})

	// Errorf logs with fmt.Sprintf-like arguments.
	Errorf(template string, args ...interface{})

	// Errorw logs with additional context using variadic key-value pairs.
	Errorw(msg string, keysAndValues ...interface{})
}

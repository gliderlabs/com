package log

type Logger interface {
	With(args ...interface{}) Logger

	DebugLogger
	InfoLogger
	WarnLogger
	ErrorLogger
}

type DebugLogger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

type InfoLogger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
}

type WarnLogger interface {
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

type ErrorLogger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

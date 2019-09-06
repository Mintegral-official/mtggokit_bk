package log

// Logger for log
type BiLogger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

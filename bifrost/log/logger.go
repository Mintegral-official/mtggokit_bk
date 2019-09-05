package log

// Logger for log
type Logger interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

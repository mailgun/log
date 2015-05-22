package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

// writerLogger is a generic type of a logger that sends messages to the underlying io.Writer.
type writerLogger struct {
	sev Severity

	w io.Writer
}

func (l *writerLogger) Infof(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityInfo, format, args...)
}

func (l *writerLogger) Warnf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityWarn, format, args...)
}

func (l *writerLogger) Errorf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityError, format, args...)
}

func (l *writerLogger) Fatalf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityFatal, format, args...)
}

func (l *writerLogger) Writer(sev Severity) io.Writer {
	// is this logger configured to log at the provided severity?
	if l.sev.Gt(sev) {
		return nil
	}
	return l.w
}

// consoleLogger is a type of writerLogger that sends messages to the standard output.
type consoleLogger struct {
	*writerLogger
}

func NewConsoleLogger(conf LogConfig) (Logger, error) {
	return &consoleLogger{&writerLogger{conf.Severity, os.Stdout}}, nil
}

func (l *writerLogger) FormatMessage(sev Severity, fileName, funcName string, lineNo int, format string, args ...interface{}) string {
	return fmt.Sprintf("%v %s %s PID:%d [%s:%d:%s] %s",
		time.Now().UTC().Format(time.StampMilli), appname, sev, pid, fileName, lineNo, funcName, fmt.Sprintf(format, args...))
}

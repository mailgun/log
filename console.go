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

func (l *writerLogger) Warningf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityWarning, format, args...)
}

func (l *writerLogger) Errorf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityError, format, args...)
}

func (l *writerLogger) Writer(sev Severity) io.Writer {
	// is this logger configured to log at the provided severity?
	if sev.Gte(l.sev) {
		return l.w
	}
	return nil
}

func (l *writerLogger) FormatMessage(sev Severity, caller *callerInfo, format string, args ...interface{}) string {
	return fmt.Sprintf("%v %s %s PID:%d [%s:%d:%s] %s\n",
		time.Now().UTC().Format(time.StampMilli), appname, sev, pid, caller.fileName, caller.lineNo, caller.funcName, fmt.Sprintf(format, args...))
}

func (l writerLogger) String() string {
	return fmt.Sprintf("writerLogger(%s)", l.sev)
}

const ConsoleLoggerName = "console"

// consoleLogger is a type of writerLogger that sends messages to the standard output.
type consoleLogger struct {
	*writerLogger
}

func NewConsoleLogger(conf LogConfig) (Logger, error) {
	sev, err := SeverityFromString(conf.Severity)
	if err != nil {
		return nil, err
	}
	return &consoleLogger{&writerLogger{sev, os.Stdout}}, nil
}

func (l consoleLogger) String() string {
	return fmt.Sprintf("consoleLogger(%s)", l.sev)
}

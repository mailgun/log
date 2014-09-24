package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Console logger is for dev mode, it prints the logs to the terminal.
// Note: don't use this logger in production.
type streamLogger struct {
	out io.Writer
}

func NewConsoleLogger(config *LogConfig) (Logger, error) {
	return &streamLogger{out: os.Stdout}, nil
}

func (l *streamLogger) Info(message string) {
	l.print(SeverityInfo, message)
}

func (l *streamLogger) Warning(message string) {
	l.print(SeverityWarn, message)
}

func (l *streamLogger) Error(message string) {
	l.print(SeverityError, message)
}

func (l *streamLogger) Fatal(message string) {
	l.print(SeverityFatal, message)
}

func (l *streamLogger) print(sev severity, message string) {
	if currentSeverity.ge(sev) {
		fmt.Fprintf(l.out, "%v %v: %v\n", sev, time.Now().UTC().Format(time.StampMilli), message)
	}
}

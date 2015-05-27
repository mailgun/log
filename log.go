package log

import (
	"errors"
	"fmt"
	"io"
)

var loggers []Logger

// Logger is a unified interface for all loggers.
type Logger interface {
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})

	Writer(Severity) io.Writer
	FormatMessage(Severity, *callerInfo, string, ...interface{}) string

	String() string
}

// Logging configuration to be passed to all loggers during initialization.
type LogConfig struct {
	Name     string
	Severity string
}

func (c LogConfig) String() string {
	return fmt.Sprintf("LogConfig(Name=%v, Severity=%v)", c.Name, c.Severity)
}

func Init(l ...Logger) {
	for _, logger := range l {
		loggers = append(loggers, logger)
	}
}

func InitWithConfig(logConfigs ...LogConfig) error {
	for _, config := range logConfigs {
		l, err := NewLogger(config)
		if err != nil {
			return err
		}
		loggers = append(loggers, l)
	}
	return nil
}

// Make a proper logger from a given configuration.
func NewLogger(config LogConfig) (Logger, error) {
	switch config.Name {
	case ConsoleLoggerName:
		return NewConsoleLogger(config)
	case SysLoggerName:
		return NewSysLogger(config)
	case UDPLoggerName:
		return NewUDPLogger(config)
	}
	return nil, errors.New(fmt.Sprintf("unknown logger: %v", config))
}

// Infof logs to the INFO log.
func Infof(format string, args ...interface{}) {
	for _, logger := range loggers {
		writeMessage(logger, 1, SeverityInfo, format, args...)
	}
}

// Warningf logs to the WARNING and INFO logs.
func Warnf(format string, args ...interface{}) {
	for _, logger := range loggers {
		writeMessage(logger, 1, SeverityWarn, format, args...)
	}
}

// Errorf logs to the ERROR, WARNING, and INFO logs.
func Errorf(format string, args ...interface{}) {
	for _, logger := range loggers {
		writeMessage(logger, 1, SeverityError, format, args...)
	}
}

// Fatalf logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
func Fatalf(format string, args ...interface{}) {
	for _, logger := range loggers {
		writeMessage(logger, 1, SeverityFatal, format, args...)
	}
}

func writeMessage(logger Logger, callDepth int, sev Severity, format string, args ...interface{}) {
	caller := getCallerInfo(callDepth + 1)
	if w := logger.Writer(sev); w != nil {
		message := logger.FormatMessage(sev, caller, format, args...)
		io.WriteString(w, message)
		if sev == SeverityFatal {
			io.WriteString(w, getStackTraces())
		}
	}
}

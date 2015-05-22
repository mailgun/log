package log

import (
	"fmt"
	"io"
	"log/syslog"
)

// Syslogger sends all your logs to syslog
// Note: the logs are going to MAIL_LOG facility
type sysLogger struct {
	sev Severity

	infoW  *syslog.Writer
	warnW  *syslog.Writer
	errorW *syslog.Writer
}

var newSyslogWriter = syslog.New // for mocking in tests

func NewSysLogger(conf LogConfig) (Logger, error) {
	infoW, err := newSyslogWriter(syslog.LOG_MAIL|syslog.LOG_INFO, appname)
	if err != nil {
		return nil, err
	}

	warnW, err := newSyslogWriter(syslog.LOG_MAIL|syslog.LOG_WARNING, appname)
	if err != nil {
		return nil, err
	}

	errorW, err := newSyslogWriter(syslog.LOG_MAIL|syslog.LOG_ERR, appname)
	if err != nil {
		return nil, err
	}

	return &sysLogger{conf.Severity, infoW, warnW, errorW}, nil
}

func (l *sysLogger) Infof(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityInfo, format, args...)
}

func (l *sysLogger) Warnf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityWarn, format, args...)
}

func (l *sysLogger) Errorf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityError, format, args...)
}

func (l *sysLogger) Fatalf(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityFatal, format, args...)
}

func (l *sysLogger) Writer(sev Severity) io.Writer {
	// is this logger configured to log at the provided severity?
	if l.sev.Gt(sev) {
		return nil
	}

	// return an appropriate writer
	switch sev {
	case SeverityInfo:
		return l.info
	case SeverityWarn:
		return l.warn
	default:
		return l.err
	}
}

func (l *sysLogger) FormatMessage(sev Severity, fileName, funcName string, lineNo int, format string, args ...interface{}) string {
	return fmt.Sprintf("%s PID:%d [%s:%d:%s] %s", sev, pid, fileName, lineNo, funcName, fmt.Sprintf(format, args...))
}

package log

import (
	"fmt"
	"io"
	"log/syslog"
)

const SysLoggerName = "syslog"

// sysLogger logs messages to rsyslog MAIL_LOG facility.
type sysLogger struct {
	sev Severity

	infoW  io.Writer
	warnW  io.Writer
	errorW io.Writer
}

func NewSysLogger(conf LogConfig) (Logger, error) {
	infoW, err := syslog.New(syslog.LOG_MAIL|syslog.LOG_INFO, appname)
	if err != nil {
		return nil, err
	}

	warnW, err := syslog.New(syslog.LOG_MAIL|syslog.LOG_WARNING, appname)
	if err != nil {
		return nil, err
	}

	errorW, err := syslog.New(syslog.LOG_MAIL|syslog.LOG_ERR, appname)
	if err != nil {
		return nil, err
	}

	sev, err := SeverityFromString(conf.Severity)
	if err != nil {
		return nil, err
	}

	return &sysLogger{sev, infoW, warnW, errorW}, nil
}

func (l *sysLogger) Info(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityInfo, format, args...)
}

func (l *sysLogger) Warning(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityWarning, format, args...)
}

func (l *sysLogger) Error(format string, args ...interface{}) {
	writeMessage(l, 1, SeverityError, format, args...)
}

func (l *sysLogger) Writer(sev Severity) io.Writer {
	// is this logger configured to log at the provided severity?
	if sev.Gte(l.sev) {
		// return an appropriate writer
		switch sev {
		case SeverityInfo:
			return l.infoW
		case SeverityWarning:
			return l.warnW
		default:
			return l.errorW
		}
	}
	return nil
}

func (l *sysLogger) FormatMessage(sev Severity, caller *callerInfo, format string, args ...interface{}) string {
	return fmt.Sprintf("%s PID:%d [%s:%d:%s] %s", sev, pid, caller.fileName, caller.lineNo, caller.funcName, fmt.Sprintf(format, args...))
}

func (l sysLogger) String() string {
	return fmt.Sprintf("sysLogger(%s)", l.sev)
}

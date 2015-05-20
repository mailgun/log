package log

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os"
	"time"
)

const (
	DefaultAddr = "127.0.0.1:55647"
)

type record struct {
	App       string   `json:"appname"`
	Host      string   `json:"hostname"`
	Sev       Severity `json:"loglevel"`
	Message   []byte   `json:"message"`
	Timestamp int64    `json:"timestamp"`
}

type udpLogger struct {
	info *udpLogWriter
	warn *udpLogWriter
	err  *udpLogWriter
}

type udpLogWriter struct {
	conn *net.UDPConn
	app  string
	host string
	sev  Severity
}

func NewUDPLogger(config *LogConfig) (Logger, error) {
	addr, err := net.ResolveUDPAddr("udp", DefaultAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &udpLogger{
		&udpLogWriter{conn, getAppName(), host, SeverityInfo},
		&udpLogWriter{conn, getAppName(), host, SeverityWarn},
		&udpLogWriter{conn, getAppName(), host, SeverityError},
	}, nil
}

func (l *udpLogger) Writer(sev Severity) io.Writer {
	switch sev {
	case SeverityInfo:
		return l.info
	case SeverityWarn:
		return l.warn
	default:
		return l.err
	}
}

func (l *udpLogger) Infof(format string, args ...interface{}) {
	infof(1, l.Writer(SeverityInfo), format, args...)
}

func (l *udpLogger) Warningf(format string, args ...interface{}) {
	warningf(1, l.Writer(SeverityWarn), format, args...)
}

func (l *udpLogger) Errorf(format string, args ...interface{}) {
	errorf(1, l.Writer(SeverityError), format, args...)
}

func (l *udpLogger) Fatalf(format string, args ...interface{}) {
	fatalf(1, l.Writer(SeverityFatal), format, args...)
}

func (l *udpLogWriter) Write(msg []byte) (int, error) {
	record := &record{l.app, l.host, l.sev, msg, time.Now().UnixNano() / 1000000}

	dump, err := json.Marshal(record)
	if err != nil {
		return 0, err
	}

	buf := bytes.Buffer{}
	buf.Write([]byte("go_udplog:\t"))
	buf.Write(dump)

	return l.conn.Write(buf.Bytes())
}

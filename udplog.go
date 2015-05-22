package log

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	DefaultHost = "127.0.0.1"
	DefaultPort = 55647

	DefaultCategory = "go_logging"
)

type udpLogRecord struct {
	AppName   string   `json:"appname"`
	HostName  string   `json:"hostname"`
	LogLevel  Severity `json:"logLevel"`
	FileName  string   `json:"filename"`
	FuncName  string   `json:"funcName"`
	LineNo    int      `json:"lineno"`
	Message   []byte   `json:"message"`
	Timestamp int64    `json:"timestamp"`
}

// udpLoggers is a type of writerLogger that sends messages in a special format to a udplog server.
type udpLogger struct {
	*writerLogger
}

func NewUDPLogger(conf LogConfig) (Logger, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", DefaultHost, DefaultPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &udpLogger{&writerLogger{conf.Severity, conn}}, nil
}

func (l *udpLogger) FormatMessage(sev Severity, fileName, funcName string, lineNo int, format string, args ...interface{}) string {
	rec := &udpLogRecord{
		appname, hostname, sev, fileName, funcName, lineNo, fmt.Sprintf(format, args...), time.Now().UnixNano() / 1000000}

	dump, err := json.Marshal(rec)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", DefaultCategory, dump)
}

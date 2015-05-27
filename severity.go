package log

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type Severity int32

const (
	SeverityInfo Severity = iota
	SeverityWarn
	SeverityError
	SeverityFatal
)

var severityName = map[Severity]string{
	SeverityInfo:  "INFO",
	SeverityWarn:  "WARN",
	SeverityError: "ERROR",
	SeverityFatal: "FATAL",
}

func (s *Severity) Get() Severity {
	return Severity(atomic.LoadInt32((*int32)(s)))
}

func (s *Severity) Set(val Severity) {
	atomic.StoreInt32((*int32)(s), int32(val))
}

func (s *Severity) Gte(val Severity) bool {
	return s.Get() >= val
}

func (s Severity) String() string {
	n, ok := severityName[s]
	if !ok {
		return "UNKNOWN SEVERITY"
	}
	return n
}

func SeverityFromString(s string) (Severity, error) {
	s = strings.ToUpper(s)
	for k, val := range severityName {
		if val == s {
			return k, nil
		}
	}
	return -1, fmt.Errorf("unsupported severity: %s", s)
}

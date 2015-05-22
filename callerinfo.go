package log

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	UnknownFile = "unknown_file"
	UnknownFunc = "unknown_func"
)

var (
	pid         = os.Getpid()
	hostname, _ = os.Hostname()
	appname     = filepath.Base(os.Args[0])
)

// callerInfo returns information about a certain log function invoker
// such as file name, function name and line number
func callerInfo(depth int) (string, string, int) {
	if pc, fileName, lineNo, ok := runtime.Caller(depth + 1); !ok {
		return UnknownFile, UnknownFunc, 0
	} else {
		slashPos := strings.LastIndex(file, "/")
		if slashPos >= 0 {
			fileName = fileName[slashPos+1:]
		}
		return fileName, runtime.FuncForPC(pc).Name(), lineNo
	}
}

// Return stack traces of all the running goroutines.
func stackTraces() string {
	trace := make([]byte, 100000)
	nbytes := runtime.Stack(trace, true)
	return string(trace[:nbytes])
}

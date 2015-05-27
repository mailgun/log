package log

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	UnknownFile = "unknown_file"
	UnknownPath = "unknown_path"
	UnknownFunc = "unknown_func"
)

var (
	pid         = os.Getpid()
	hostname, _ = os.Hostname()
	appname     = filepath.Base(os.Args[0])
)

type callerInfo struct {
	fileName string
	filePath string
	funcName string
	lineNo   int
}

// getCallerInfo returns information about a certain log function invoker
// such as file name, function name and line number
func getCallerInfo(depth int) *callerInfo {
	if pc, filePath, lineNo, ok := runtime.Caller(depth + 1); !ok {
		return &callerInfo{UnknownFile, UnknownPath, UnknownFunc, 0}
	} else {
		var fileName string
		slashPos := strings.LastIndex(filePath, "/")
		if slashPos >= 0 {
			fileName = filePath[slashPos+1:]
		}
		return &callerInfo{fileName, filePath, runtime.FuncForPC(pc).Name(), lineNo}
	}
}

package log

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestModel(t *testing.T) { TestingT(t) }

type LogSuite struct{}

var _ = Suite(&LogSuite{})

func (s *LogSuite) SetUpTest(c *C) {
	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return 0, "", 0, false
	}
	// mock exit function
	exit = func() {}
}

func (s *LogSuite) SetUpSuite(c *C) {
	consoleConfig := &LogConfig{Name: "console"}
	syslogConfig := &LogConfig{Name: "syslog"}
	Init([]*LogConfig{consoleConfig, syslogConfig})
}

func (s *LogSuite) TestInit(c *C) {
	consoleConfig := &LogConfig{Name: "console"}
	syslogConfig := &LogConfig{Name: "syslog"}

	err := Init([]*LogConfig{consoleConfig, syslogConfig})
	c.Assert(err, IsNil)
}

func (s *LogSuite) TestInitError(c *C) {
	unknownConfig := &LogConfig{Name: "unknown"}
	err := Init([]*LogConfig{unknownConfig})
	c.Assert(err, NotNil)
}

func (s *LogSuite) TestInfof(c *C) {
	Infof("test message, %v", "info")
}

func (s *LogSuite) TestWarningf(c *C) {
	Warningf("test message, %v", "warning")
}

func (s *LogSuite) TestErrorf(c *C) {
	Errorf("test message, %v", "error")
}

func (s *LogSuite) TestFatalf(c *C) {

	Fatalf("test message, %v", "fatal")
}

func (s *LogSuite) TestCallerInfoError(c *C) {
	file, line := callerInfo()
	c.Assert(file, Equals, "unknown")
	c.Assert(line, Equals, 0)
}

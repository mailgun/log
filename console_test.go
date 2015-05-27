package log

import (
	"bytes"
	"os"
	"strings"

	. "gopkg.in/check.v1"
)

type WriterLoggerSuite struct {
	w *bytes.Buffer
	l *writerLogger
}

var _ = Suite(&WriterLoggerSuite{})

func (s *WriterLoggerSuite) SetUpTest(c *C) {
	s.w = &bytes.Buffer{}
	s.l = &writerLogger{SeverityInfo, s.w}
}

func (s *WriterLoggerSuite) output() string {
	return s.w.String()
}

func (s *WriterLoggerSuite) TestInfof(c *C) {
	s.l.Infof("log message")
	c.Assert(s.output(), Matches, ".*INFO.*log message.*\n")
}

func (s *WriterLoggerSuite) TestWarnf(c *C) {
	s.l.Warnf("log message")
	c.Assert(s.output(), Matches, ".*WARN.*log message.*\n")
}

func (s *WriterLoggerSuite) TestErrorf(c *C) {
	s.l.Errorf("log message")
	c.Assert(s.output(), Matches, ".*ERROR.*log message.*\n")
}

func (s *WriterLoggerSuite) TestFatalf(c *C) {
	s.l.Fatalf("log message")
	c.Assert(strings.Split(s.output(), "\n")[0], Matches, ".*FATAL.*log message.*")
}

func (s *WriterLoggerSuite) TestSeverity(c *C) {
	// create an error logger
	l := &writerLogger{SeverityError, s.w}

	// it should not log anything below ERROR
	l.Infof("log message")
	c.Assert(s.output(), Equals, "")

	l.Warnf("log message")
	c.Assert(s.output(), Equals, "")

	l.Errorf("log message")
	c.Assert(s.output(), Matches, ".*ERROR.*log message.*\n")
}

type ConsoleLoggerSuite struct {
}

var _ = Suite(&ConsoleLoggerSuite{})

func (s *ConsoleLoggerSuite) TestNewConsoleLogger(c *C) {
	l, err := NewConsoleLogger(LogConfig{ConsoleLoggerName, "info"})
	c.Assert(err, IsNil)
	c.Assert(l, NotNil)

	console := l.(*consoleLogger)
	c.Assert(console.sev, Equals, SeverityInfo)
	c.Assert(console.w, Equals, os.Stdout)
}

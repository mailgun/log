package log

import (
	"bytes"
	"os"

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

func (s *WriterLoggerSuite) TestInfo(c *C) {
	s.l.Info("log message")
	c.Assert(s.output(), Matches, ".*INFO.*log message.*\n")
}

func (s *WriterLoggerSuite) TestWarning(c *C) {
	s.l.Warning("log message")
	c.Assert(s.output(), Matches, ".*WARN.*log message.*\n")
}

func (s *WriterLoggerSuite) TestError(c *C) {
	s.l.Error("log message")
	c.Assert(s.output(), Matches, ".*ERROR.*log message.*\n")
}

func (s *WriterLoggerSuite) TestSeverity(c *C) {
	// create an error logger
	l := &writerLogger{SeverityError, s.w}

	// it should not log anything below ERROR
	l.Info("log message")
	c.Assert(s.output(), Equals, "")

	l.Warning("log message")
	c.Assert(s.output(), Equals, "")

	l.Error("log message")
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

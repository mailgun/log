package log

import (
	"bytes"
	"strings"

	. "gopkg.in/check.v1"
)

type ConsoleLogSuite struct {
	logger Logger
	out    *bytes.Buffer
}

var _ = Suite(&ConsoleLogSuite{})

func (s *ConsoleLogSuite) SetUpTest(c *C) {
	SetSeverity(SeverityInfo)
	s.out = &bytes.Buffer{}
	s.logger = &streamLogger{out: s.out}
}

func (s *ConsoleLogSuite) output() string {
	return strings.TrimSpace(s.out.String())
}

func (s *ConsoleLogSuite) TestNewConsoleLogger(c *C) {
	config := &LogConfig{Name: "testNew"}
	logger, err := NewConsoleLogger(config)
	c.Assert(logger, NotNil)
	c.Assert(err, IsNil)
}

func (s *ConsoleLogSuite) TestInfo(c *C) {
	s.logger.Info("test message")
	c.Assert(s.output(), Matches, "INFO.*test message.*")
}

func (s *ConsoleLogSuite) TestWarning(c *C) {
	s.logger.Warning("test message")
	c.Assert(s.output(), Matches, "WARNING.*test message.*")
}

func (s *ConsoleLogSuite) TestError(c *C) {
	s.logger.Error("test message")
	c.Assert(s.output(), Matches, "ERROR.*test message.*")
}

func (s *ConsoleLogSuite) TestFatal(c *C) {
	s.logger.Fatal("test message")
	c.Assert(s.output(), Matches, "FATAL.*test message.*")
}

func (s *ConsoleLogSuite) TestUpperLevel(c *C) {
	SetSeverity(SeverityError)
	s.logger.Info("info message")
	s.logger.Error("error message")
	c.Assert(s.output(), Matches, "ERROR.*error message.*")
}

func (s *ConsoleLogSuite) TestUpdateLevel(c *C) {
	SetSeverity(SeverityError)
	s.logger.Info("info message")
	c.Assert(s.output(), Equals, "")

	SetSeverity(SeverityInfo)
	s.logger.Info("info message")
	c.Assert(s.output(), Matches, "INFO.*info message.*")
}

package log

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type SysLoggerSuite struct {
	l *sysLogger

	infoB  *bytes.Buffer
	warnB  *bytes.Buffer
	errorB *bytes.Buffer
}

var _ = Suite(&SysLoggerSuite{})

func (s *SysLoggerSuite) SetUpTest(c *C) {
	s.infoB = &bytes.Buffer{}
	s.warnB = &bytes.Buffer{}
	s.errorB = &bytes.Buffer{}
	s.l = &sysLogger{SeverityInfo, s.infoB, s.warnB, s.errorB}
}

func (s *SysLoggerSuite) TestInfof(c *C) {
	s.l.Infof("info message")
	c.Assert(s.infoB.String(), Matches, ".*INFO.*info message.*")
	c.Assert(s.warnB.String(), Equals, "")
	c.Assert(s.errorB.String(), Equals, "")
}

func (s *SysLoggerSuite) TestWarningf(c *C) {
	s.l.Warningf("warn message")
	c.Assert(s.infoB.String(), Equals, "")
	c.Assert(s.warnB.String(), Matches, ".*WARN.*warn message.*")
	c.Assert(s.errorB.String(), Equals, "")
}

func (s *SysLoggerSuite) TestErrorf(c *C) {
	s.l.Errorf("error message")
	c.Assert(s.infoB.String(), Equals, "")
	c.Assert(s.warnB.String(), Equals, "")
	c.Assert(s.errorB.String(), Matches, ".*ERROR.*error message.*")
}

func (s *SysLoggerSuite) TestSeverity(c *C) {
	// create an error logger
	l := &sysLogger{SeverityError, s.infoB, s.warnB, s.errorB}

	// it should not log anything below ERROR
	l.Infof("info message")
	c.Assert(s.infoB.String(), Equals, "")

	l.Warningf("warn message")
	c.Assert(s.warnB.String(), Equals, "")

	l.Errorf("error message")
	c.Assert(s.errorB.String(), Matches, ".*ERROR.*error message.*")
}

func (s *SysLoggerSuite) TestNewSysLogger(c *C) {
	l, err := NewSysLogger(LogConfig{SysLoggerName, "info"})
	c.Assert(err, IsNil)
	c.Assert(l, NotNil)

	syslog := l.(*sysLogger)
	c.Assert(syslog.sev, Equals, SeverityInfo)
	c.Assert(syslog.infoW, NotNil)
	c.Assert(syslog.warnW, NotNil)
	c.Assert(syslog.errorW, NotNil)
}

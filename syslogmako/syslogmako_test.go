package syslogmako

import (
	"testing"

	. "launchpad.net/gocheck"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type SyslogMakoTestSuite struct {
}

var _ = Suite(&SyslogMakoTestSuite{})

func (s *SyslogMakoTestSuite) TestParser_Valid(c *C) {
	c.Assert(nil, Equals, nil)
}

package mako

import (
	"testing"

	. "launchpad.net/gocheck"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type MakoTestSuite struct {
}

var _ = Suite(&MakoTestSuite{})

func (s *MakoTestSuite) TestParser_Valid(c *C) {
	c.Assert(nil, Equals, nil)
}

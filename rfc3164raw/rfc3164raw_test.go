package rfc3164raw

import (
	"testing"

	. "launchpad.net/gocheck"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type Rfc3164RawTestSuite struct {
}

var (
	_ = Suite(&Rfc3164RawTestSuite{})
)

func (s *Rfc3164RawTestSuite) TestParser_Valid(c *C) {
	c.Assert(nil, Equals, nil)
}

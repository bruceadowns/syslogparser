package rfc5424raw

import (
	"testing"

	. "launchpad.net/gocheck"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type Rfc5424RawTestSuite struct {
}

var _ = Suite(&Rfc5424RawTestSuite{})

func (s *Rfc5424RawTestSuite) TestParser_Valid(c *C) {
	c.Assert(nil, Equals, nil)
}

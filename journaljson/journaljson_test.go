package journaljson

import (
	"testing"

	. "launchpad.net/gocheck"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type JournalJSONTestSuite struct {
}

var _ = Suite(&JournalJSONTestSuite{})

func (s *JournalJSONTestSuite) TestParser_Valid(c *C) {
	c.Assert(nil, Equals, nil)
}

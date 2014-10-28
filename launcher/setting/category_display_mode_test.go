package setting

import (
	"fmt"
	C "launchpad.net/gocheck"
)

type CategoryDisplayModeTestSuite struct {
}

var _ = C.Suite(CategoryDisplayModeTestSuite{})

func (sts CategoryDisplayModeTestSuite) TestCategoryDisplayMode(c *C.C) {
	c.Assert(fmt.Sprint(CategoryDisplayModeUnknown), C.Equals, "unknown category display mode")
	c.Assert(fmt.Sprint(CategoryDisplayModeText), C.Equals, "display text mode")
	c.Assert(fmt.Sprint(CategoryDisplayModeIcon), C.Equals, "display icon mode")
}

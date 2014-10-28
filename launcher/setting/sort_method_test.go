package setting

import (
	"fmt"
	C "launchpad.net/gocheck"
)

type SortMethodTestSuite struct {
}

var _ = C.Suite(&SortMethodTestSuite{})

func (sts *SortMethodTestSuite) TestSortMethod(c *C.C) {
	c.Assert(fmt.Sprint(SortMethodUnknown), C.Equals, "unknown sort method")
	c.Assert(fmt.Sprint(SortMethodByName), C.Equals, "sort by name")
	c.Assert(fmt.Sprint(SortMethodByCategory), C.Equals, "sort by category")
	c.Assert(fmt.Sprint(SortMethodByTimeInstalled), C.Equals, "sort by time installed")
	c.Assert(fmt.Sprint(SortMethodByFrequency), C.Equals, "sort by frequency")
}

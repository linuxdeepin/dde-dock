/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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

/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package search

import (
	C "launchpad.net/gocheck"
	"regexp"
)

type MatcherTestSuite struct {
}

var _ = C.Suite(&MatcherTestSuite{})

func (*MatcherTestSuite) TestMatcher(c *C.C) {
	// TODO: test them
	getMatchers("firefox")
	getMatchers("深度")
	getMatchers("f")
	getMatchers(regexp.QuoteMeta("f\\"))
}

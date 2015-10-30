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

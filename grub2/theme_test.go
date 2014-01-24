package main

import (
	. "launchpad.net/gocheck"
)

var theme *Theme

func init() {
	theme = &Theme{}
	Suite(theme)
}

// TODO
func (theme *Theme) TestGetterAndSetter(c *C) {
}

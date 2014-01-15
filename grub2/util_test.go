package main

import (
	. "launchpad.net/gocheck"
)

type Util struct{}

var _ = Suite(&Util{})

func (u *Util) TestQuoteString(c *C) {
	var tests = []struct {
		s, want string
	}{
		{`abc`, `"abc"`},
		{``, `""`},
	}
	for _, t := range tests {
		c.Check(quoteString(t.s), Equals, t.want)
	}
}

func (u *Util) TestUnquoteString(c *C) {
	var tests = []struct {
		s, want string
	}{
		{`"abc"`, `abc`},
		{`'abc'`, `abc`},
		{`"abc`, `"abc`},
		{`'abc`, `'abc`},
		{`abc`, `abc`},
		{`  "abc"`, `  "abc"`},
		{``, ``},
	}
	for _, t := range tests {
		c.Check(unquoteString(t.s), Equals, t.want)
	}
}

package grub2

import (
	. "launchpad.net/gocheck"
)

func (*GrubTester) TestQuoteString(c *C) {
	var tests = []struct {
		s, want string
	}{
		{`abc`, `"abc"`},
		{`/abc/`, `"/abc/"`},
		{`abc def`, `"abc def"`},
		{`abc\def`, `"abc\\def"`},
		{``, `""`},
	}
	for _, t := range tests {
		c.Check(quoteString(t.s), Equals, t.want)
	}
}

func (*GrubTester) TestUnquoteString(c *C) {
	var tests = []struct {
		s, want string
	}{
		{`"abc"`, `abc`},
		{`'abc'`, `abc`},
		{`"abc`, `"abc`},
		{`'abc`, `'abc`},
		{`abc`, `abc`},
		{`  "abc"`, `  "abc"`},
		{`"abc def"`, `abc def`},
		{`"abc\\def"`, `abc\def`},
		{`"/abc/"`, `/abc/`},
		{``, ``},
	}
	for _, t := range tests {
		c.Check(unquoteString(t.s), Equals, t.want)
	}
}

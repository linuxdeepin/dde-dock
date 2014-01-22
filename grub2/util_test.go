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
		{`/abc/`, `"/abc/"`},
		{`abc def`, `"abc def"`},
		{`abc\def`, `"abc\\def"`},
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
		{`"abc def"`, `abc def`},
		{`"abc\\def"`, `abc\def`},
		{`"/abc/"`, `/abc/`},
		{``, ``},
	}
	for _, t := range tests {
		c.Check(unquoteString(t.s), Equals, t.want)
	}
}

func (u *Util) TestGetPathLevel(c *C) {
	var tests = []struct {
		s    string
		want int
	}{
		{"/a/b/c", 3},
		{"a/b/c/", 3},
		{"/a/b/c/", 3},
		{"a/b/c", 3},
		{"/a/b/c/file", 4},
		{"/", 1},
		{"/file", 1},
		{".", 0},
		{"./", 0},
		{"./file", 1},
		{"./././file", 1},
		{"file", 1},
		{"", 0},
	}
	for _, t := range tests {
		c.Check(getPathLevel(t.s), Equals, t.want)
	}
}

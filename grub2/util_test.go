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

func (u *Util) TestGetImgClipSizeByResolution(c *C) {
	var tests = []struct {
		sw, sh, iw, ih, wantw, wanth int32
	}{
		{1024, 768, 1024, 768, 1024, 768},
		{1024, 768, 1920, 1080, 1024, 768},
	}
	for _, t := range tests {
		w, h := getImgClipSizeByResolution(t.sw, t.sh, t.iw, t.ih)
		c.Check(w, Equals, t.wantw)
		c.Check(h, Equals, t.wanth)
	}
	var iw, ih int32 = 1920, 1080
	for sw := 1; sw < 3000; sw += 5 {
		for sh := 1; sh < 3000; sh += 5 {
			w, h := getImgClipSizeByResolution(int32(sw), int32(sh), 1920, 1080)
			if w > iw || h > ih {
				c.Fatalf("sw=%d, sh=%d, iw=%d, ih=%d, w=%d, h=%d", sw, sh, iw, ih, w, h)
			}
		}
	}
}

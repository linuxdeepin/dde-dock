package grub2

import (
	C "launchpad.net/gocheck"
)

func (*GrubTester) TestQuoteString(c *C.C) {
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
		c.Check(quoteString(t.s), C.Equals, t.want)
	}
}

func (*GrubTester) TestUnquoteString(c *C.C) {
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
		c.Check(unquoteString(t.s), C.Equals, t.want)
	}
}

func (*GrubTester) TestConvertToSimpleEntry(c *C.C) {
	data := []struct {
		v, r string
	}{
		{"", ""},
		{">", ""},
		{"0", "0"},
		{"0>1", "0"},
		{"1>1>2", "1"},
		{"Parent Title>Child Title", "Parent Title"},
	}
	for _, d := range data {
		c.Check(convertToSimpleEntry(d.v), C.Equals, d.r)
	}
}

func (*GrubTester) TestParseGfxmode(c *C.C) {
	sw, sh := getPrimaryScreenBestResolution()
	data := []struct {
		v    string
		w, h uint16
	}{
		{"auto", sw, sh},
		{"auto,800x600", sw, sh},
		{"1024x768", 1024, 768},
		{"1024x768x24", 1024, 768},
		{"1024x768,800x600,auto", 1024, 768},
		{"1024x768;800x600;auto", 1024, 768},
		{"1024x768x24,800x600,auto", 1024, 768},
	}
	for _, d := range data {
		w, h, _ := doParseGfxmode(d.v)
		c.Check(w, C.Equals, d.w)
		c.Check(h, C.Equals, d.h)
	}

	// test wrong format
	_, _, err := doParseGfxmode("")
	c.Check(err, C.NotNil)
	_, _, err = doParseGfxmode("1024")
	c.Check(err, C.NotNil)
	_, _, err = doParseGfxmode("1024x")
	c.Check(err, C.NotNil)
	_, _, err = doParseGfxmode("autox24")
	c.Check(err, C.NotNil)
}

package main

import (
	"testing"
)

func TestQuoteString(t *testing.T) {
	var tests = []struct {
		s, want string
	}{
		{`abc`, `"abc"`},
		{``, `""`},
	}
	for _, c := range tests {
		got := quoteString(c.s)
		if got != c.want {
			t.Errorf("quoteString(%q) == %q, want %q", c.s, got, c.want)
		}
	}
}

func TestUnquoteString(t *testing.T) {
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
	for _, c := range tests {
		got := unquoteString(c.s)
		if got != c.want {
			t.Errorf("unquoteString(%q) == %q, want %q", c.s, got, c.want)
		}
	}
}
